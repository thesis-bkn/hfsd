import torch
import datetime
import os
from PIL import Image
from sqlalchemy.engine.base import Connection
import io

from src.database.models import ModelStatus, Task
from src.database.query import Querier, UpdateSampleTasksParams, SaveSampleAssetParams
from src.handlers import utils
from src.handlers.pipeline import pipeline_with_logprob_inpaint
from src.s3 import ImageUploader
from src.handlers import consts


class SampleHander:
    def __init__(self, conn: Connection, uploader: ImageUploader) -> None:
        self.conn = conn
        self.uploader = uploader
        self.querier = Querier(self.conn)

    def handle(self, task: Task):
        pipe = utils.prepare_pipe()

        source_model = self.querier.get_model(id=task.source_model_id)  # pyright: ignore
        if source_model is None:
            print("model not found")
            exit(1)

        if "base" not in source_model.name and source_model.ckpt is not None:
            pipe.unet.load_attn_procs(
                torch.load(io.BytesIO(source_model.ckpt.tobytes())), weights_only=True
            )

        prompt, neg_prompt = utils.get_prompt(source_model.domain)  # pyright: ignore
        neg_prompt_embed = pipe.text_encoder(
            pipe.tokenizer(
                [neg_prompt],
                return_tensors="pt",
                padding="max_length",
                truncation=True,
                max_length=pipe.tokenizer.model_max_length,
            ).input_ids.to(pipe.device)
        )[0]

        pipe.unet.eval()
        sample_neg_prompt_embeds = neg_prompt_embed.repeat(consts.BATCH_SIZE, 1, 1)

        # Sampling
        random_base_assets = self.querier.get_random_base_assets_by_domain(
            domain=source_model.domain, limit=consts.BATCH_SIZE
        )
        images = []
        masks = []
        for base_asset in random_base_assets:
            images.append(mem_to_pil(base_asset.image))
            masks.append(mem_to_pil(base_asset.mask))

        prompts_embeds = list(
            map(
                lambda p: pipe.text_encoder(
                    pipe.tokenizer(
                        p,
                        return_tensors="pt",
                        padding="max_length",
                        truncation=True,
                        max_length=pipe.tokenizer.model_max_length,
                    ).input_ids.to(pipe.device)
                )[0],
                [[prompt] * consts.BATCH_SIZE] * consts.NUM_PER_PROMPT,
            )
        )

        sample_results = list(
            map(
                lambda p_embed: sample(
                    pipe,
                    p_embed,
                    sample_neg_prompt_embeds,
                    images,
                    masks,
                ),
                prompts_embeds,
            )
        )

        post_images = []
        post_latents = []
        post_masks = []
        post_mask_latents = []
        for result in sample_results:
            post_images.append(result[0])
            post_latents.append(result[1])
            post_masks.append(result[2])
            post_mask_latents.append(result[3])

        # mask_latents = torch.stack(post_mask_latents, dim=1)
        image_torchs = torch.stack(post_images, dim=1)
        latents = torch.stack(post_latents, dim=1)
        prompt_embeds = torch.stack(prompts_embeds, dim=1)
        next_latents = latents[:, :, 1:]
        timesteps = pipe.scheduler.timesteps.repeat(
            consts.BATCH_SIZE, 1
        )  # (batch_size, num_steps)

        image_torchs_b = torch_to_bytes(image_torchs)
        latents_b = torch_to_bytes(latents)
        prompt_embeds_b = torch_to_bytes(prompt_embeds)
        next_latents_b = torch_to_bytes(next_latents)
        timesteps_b = torch_to_bytes(timesteps)

        print("[sample] latents shape: ", latents.shape)
        self.querier.update_sample_tasks(
            UpdateSampleTasksParams(
                id=task.id,
                latents=memoryview(latents_b),
                prompt_embeds=memoryview(prompt_embeds_b),
                next_latents=memoryview(next_latents_b),
                timesteps=memoryview(timesteps_b),
                image_torchs=memoryview(image_torchs_b),
            )
        )

        for order, image in enumerate(post_images):
            for k in range(consts.BATCH_SIZE):
                pil = Image.fromarray(
                    (image[k].cpu().numpy().transpose(1, 2, 0) * 255).astype("uint8"),
                    "RGB",
                )
                img_byte_arr = io.BytesIO()
                pil.save(img_byte_arr, format="PNG")
                img_bytes = img_byte_arr.getvalue()
                image_url = os.path.join(
                    "sample",
                    "task_{}_group_{}_order_{}".format(
                        task.id, k, order + k * consts.NUM_PER_PROMPT
                    ),
                )

                self.uploader.upload_image(img_bytes, image_url)
                self.querier.save_sample_asset(
                    SaveSampleAssetParams(
                        task_id=task.id,
                        order=order + k * consts.NUM_PER_PROMPT,
                        group=k,
                        image=memoryview(img_byte_arr.getvalue()),
                        image_url=image_url,
                        prompt=prompt,
                    )
                )

        self.querier.update_task_status(
            id=task.id,
            handled_at=None,
            finished_at=datetime.datetime.now(datetime.UTC),
        )

        self.querier.update_model_status(
            id=task.output_model_id, # pyright: ignore
            status=ModelStatus.RATING
        )
        self.conn.commit()


def sample(pipe, prompt_embed, neg_prompt_embed, input_images, input_masks):
    images, _, latents, _, mask, masklatents = pipeline_with_logprob_inpaint(
        pipe,
        image=input_images,
        mask_image=input_masks,
        prompt_embeds=prompt_embed,
        negative_prompt_embeds=neg_prompt_embed,
        num_inference_steps=consts.NUM_STEPS,
        guidance_scale=consts.GUIDANCE_SCALE,
        eta=consts.ETA,
        output_type="pt",
    )
    masklatents = torch.stack(masklatents, dim=1)
    images = images.cpu().detach()  # pyright: ignore
    latents = torch.stack(latents, dim=1).cpu().detach()  # pyright: ignore

    return images, latents, mask, masklatents


def mem_to_pil(x: memoryview):
    image = Image.open(io.BytesIO(x.tobytes()))
    return image


def torch_to_bytes(t):
    buff = io.BytesIO()
    torch.save(t, buff)
    return buff.getvalue()
