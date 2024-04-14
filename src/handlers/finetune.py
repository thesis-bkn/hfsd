import copy
import io
from itertools import combinations

import numpy as np
import torch
import tree
from sqlalchemy.engine.base import Connection
from tqdm import tqdm
from typing_extensions import List

from src.database.models import Asset, Task
from src.database.query import Querier
from src.handlers import utils
from src.handlers.ddim import ddim_step_with_logprob
from src.s3 import ImageUploader

NUMSTEPS = 2
BETA = 0.1
EPS = 0.1
GRADIENT_ACCUMULATION_STEPS = 1
TIMESTEP_FRACTION = 1.0
NUM_STEPS = 20
BATCH_SIZE = 10
GUIDANCE_SCALE = 5.0
ETA = 1.0
NUM_PER_PROMPT = 7


class FinetuneHandler:
    def __init__(self, conn: Connection, uploader: ImageUploader) -> None:
        self.conn = conn
        self.uploader = uploader
        self.querier = Querier(self.conn)

    def handle(self, task: Task):
        num_train_timesteps = int(NUM_STEPS * TIMESTEP_FRACTION)
        pipe = utils.prepare_pipe()
        optimizer = utils.prepare_optimizer(pipe)

        timesteps = torch.load(io.BytesIO(task.timesteps.tobytes()))  # pyright: ignore
        prompt_embeds = torch.load(io.BytesIO(task.prompt_embeds.tobytes()))  # pyright: ignore
        latents = torch.load(io.BytesIO(task.prompt_embeds.tobytes()))  # pyright: ignore
        next_latents = torch.load(io.BytesIO(task.next_latents.tobytes()))  # pyright: ignore
        image_torchs = torch.load(io.BytesIO(task.image_torchs.tobytes()))  # pyright: ignore
        samples = {
            "timesteps": timesteps,
            "prompt_embeds": prompt_embeds,
            "latents": latents,
            "next_latents": next_latents,
            "image_torchs": image_torchs,
        }

        source_model = self.querier.get_model(id=task.source_model_id)  # pyright: ignore
        if source_model is None:
            print("model not found")
            exit(1)

        _, neg_prompt = utils.get_prompt(source_model.domain)  # pyright: ignore
        neg_prompt_embed = pipe.text_encoder(
            pipe.tokenizer(
                [neg_prompt],
                return_tensors="pt",
                padding="max_length",
                truncation=True,
                max_length=pipe.tokenizer.model_max_length,
            ).input_ids.to(pipe.device)
        )[0]

        train_neg_prompt_embed = neg_prompt_embed.repeat(BATCH_SIZE, 1, 1)

        pipe.unet.eval()

        # load lora model
        if "base" not in source_model.name and source_model.ckpt is not None:
            pipe.unet.load_attn_procs(
                torch.load(io.BytesIO(source_model.ckpt.tobytes())), weights_only=True
            )

        pipe.scheduler.timesteps = np.load(io.BytesIO(task.timesteps.tobytes()))  # pyright: ignore
        pipe.scheduler.set_timesteps(NUM_STEPS, device=pipe.device)

        ref = copy.deepcopy(pipe.unet)
        for param in ref.parameters():
            param.requires_grad = False

        assets = list(self.querier.list_asset_by_task(task_id=task.id))
        grouped_assets = group_assets(assets)
        hfs = get_hfs(grouped_assets=grouped_assets)
        samples["human_prefer"] = hfs

        num_timesteps = timesteps.shape[1]

        assert num_timesteps == NUM_STEPS

        # Training
        total_batch_size = hfs.shape[0]
        combinations_list = list(combinations(range(7), 2))
        perm = torch.randperm(hfs.shape[0])
        samples = {k: v[perm] for k, v in samples.items()}
        # shuffle along time dimension independently for each sample
        perms = torch.stack(
            [torch.randperm(num_timesteps) for _ in range(total_batch_size)]
        )

        for key in ["latents", "next_latents"]:
            tmp = samples[key].permute(0, 2, 3, 4, 5, 1)[
                torch.arange(total_batch_size)[:, None], perms
            ]
            samples[key] = tmp.permute(0, 5, 1, 2, 3, 4)
            del tmp
        samples["timesteps"] = (
            samples["timesteps"][torch.arange(total_batch_size)[:, None], perms]
            .unsqueeze(1)
            .repeat(1, 7, 1)
        )
        pipe.unet.train()
        for each_combination in combinations_list:
            sample_0 = tree.map_structure(
                lambda x: x[0, each_combination[0]].to(pipe.device),
                samples,
            )
            sample_1 = tree.map_structure(
                lambda x: x[0, each_combination[1]].to(pipe.device),
                samples,
            )
            if torch.all(sample_0["human_prefer"] == sample_1["human_prefer"]):  # pyright: ignore
                continue
            # compute which image is better
            compare_sample0 = (
                sample_0["human_prefer"] > sample_1["human_prefer"]  # pyright: ignore
            ).int() * 2 - 1
            compare_sample1 = (
                sample_1["human_prefer"] > sample_0["human_prefer"]  # pyright: ignore
            ).int() * 2 - 1
            equal_mask = sample_0["human_prefer"] == sample_1["human_prefer"]  # pyright: ignore
            compare_sample0[equal_mask] = 0
            compare_sample1[equal_mask] = 0
            human_prefer = torch.stack([compare_sample0, compare_sample1], dim=1)

            # concat negative prompts to sample prompts to avoid two forward passes
            embeds_0 = torch.cat(
                [train_neg_prompt_embed, sample_0["prompt_embeds"]]  # pyright: ignore
            )
            embeds_1 = torch.cat(
                [train_neg_prompt_embed, sample_1["prompt_embeds"]]  # pyright: ignore
            )

            for j in tqdm(
                range(num_train_timesteps),
                desc="Timestep",
                position=3,
                leave=False,
            ):
                noise_pred_0 = pipe.unet(
                    torch.cat([sample_0["latents"][:, j]] * 2),  # pyright: ignore
                    torch.cat([sample_0["timesteps"][:, j]] * 2),  # pyright: ignore
                    embeds_0,
                ).sample
                (
                    noise_pred_uncond_0,
                    noise_pred_text_0,
                ) = noise_pred_0.chunk(2)
                noise_pred_0 = noise_pred_uncond_0 + GUIDANCE_SCALE * (
                    noise_pred_text_0 - noise_pred_uncond_0
                )

                noise_ref_pred_0 = ref(
                    torch.cat([sample_0["latents"][:, j]] * 2),  # pyright: ignore
                    torch.cat([sample_0["timesteps"][:, j]] * 2),  # pyright: ignore
                    embeds_0,
                ).sample
                (
                    noise_ref_pred_uncond_0,
                    noise_ref_pred_text_0,
                ) = noise_ref_pred_0.chunk(2)
                noise_ref_pred_0 = noise_ref_pred_uncond_0 + GUIDANCE_SCALE * (
                    noise_ref_pred_text_0 - noise_ref_pred_uncond_0
                )

                noise_pred_1 = pipe.unet(
                    torch.cat([sample_1["latents"][:, j]] * 2),  # pyright: ignore
                    torch.cat([sample_1["timesteps"][:, j]] * 2),  # pyright: ignore
                    embeds_1,
                ).sample
                (
                    noise_pred_uncond_1,
                    noise_pred_text_1,
                ) = noise_pred_1.chunk(2)
                noise_pred_1 = noise_pred_uncond_1 + GUIDANCE_SCALE * (
                    noise_pred_text_1 - noise_pred_uncond_1
                )

                noise_ref_pred_1 = ref(
                    torch.cat([sample_1["latents"][:, j]] * 2),  # pyright: ignore
                    torch.cat([sample_1["timesteps"][:, j]] * 2),  # pyright: ignore
                    embeds_1,
                ).sample
                (
                    noise_ref_pred_uncond_1,
                    noise_ref_pred_text_1,
                ) = noise_ref_pred_1.chunk(2)
                noise_ref_pred_1 = noise_ref_pred_uncond_1 + GUIDANCE_SCALE * (
                    noise_ref_pred_text_1 - noise_ref_pred_uncond_1
                )

                # compute the log prob of next_latents given latents under the current model
                _, total_prob_0 = ddim_step_with_logprob(
                    pipe.scheduler,
                    noise_pred_0,
                    sample_0["timesteps"][:, j],  # pyright: ignore
                    sample_0["latents"][:, j],  # pyright: ignore
                    eta=ETA,
                    prev_sample=sample_0["next_latents"][:, j],  # pyright: ignore
                )
                _, total_ref_prob_0 = ddim_step_with_logprob(
                    pipe.scheduler,  # pyright: ignore
                    noise_ref_pred_0,
                    sample_0["timesteps"][:, j],  # pyright: ignore
                    sample_0["latents"][:, j],  # pyright: ignore
                    eta=ETA,
                    prev_sample=sample_0["next_latents"][:, j],  # pyright: ignore
                )
                _, total_prob_1 = ddim_step_with_logprob(
                    pipe.scheduler,
                    noise_pred_1,
                    sample_1["timesteps"][:, j],  # pyright: ignore
                    sample_1["latents"][:, j],  # pyright: ignore
                    eta=ETA,
                    prev_sample=sample_1["next_latents"][:, j],  # pyright: ignore
                )
                _, total_ref_prob_1 = ddim_step_with_logprob(
                    pipe.scheduler,
                    noise_ref_pred_1,
                    sample_1["timesteps"][:, j],  # pyright: ignore
                    sample_1["latents"][:, j],  # pyright: ignore
                    eta=ETA,
                    prev_sample=sample_1["next_latents"][:, j],  # pyright: ignore
                )
        # clip the probs of the pre-trained model and this model
        ratio_0 = torch.clamp(
            torch.exp(total_prob_0 - total_ref_prob_0),  # pyright: ignore
            1 - EPS,
            1 + EPS,
        )
        ratio_1 = torch.clamp(
            torch.exp(total_prob_1 - total_ref_prob_1),  # pyright: ignore
            1 - EPS,
            1 + EPS,
        )
        loss = -torch.log(
            torch.sigmoid(
                BETA * (torch.log(ratio_0)) * hfs[:, 0]
                + BETA * (torch.log(ratio_1)) * hfs[:, 1]
            )
        ).mean()

        # backward pass
        pipe.backward(loss)
        optimizer.step()
        optimizer.zero_grad()

        # save model here


def group_assets(assets: list[Asset]) -> List[List[Asset]]:
    d = {}
    for asset in assets:
        if asset.group is None:
            exit(1)
        if asset.group not in d:
            d[asset.group] = []
        d[asset.group].append(asset)

    return list(d.values())


def get_hfs(grouped_assets: List[List[Asset]]):
    return np.array(
        tree.map_structure(lambda x: x.pref, grouped_assets),
        dtype=np.float64,
    )
