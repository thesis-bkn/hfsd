import torch
from diffusers.pipelines.stable_diffusion.pipeline_stable_diffusion_inpaint import (
    StableDiffusionInpaintPipeline,
)
from diffusers.schedulers.scheduling_ddim import DDIMScheduler
from sqlalchemy.engine.base import Connection

from src.database.models import Task
from src.database.query import Querier
from src.handlers import utils
from src.handlers.pipeline import pipeline_with_logprob_inpaint
from src.s3 import ImageUploader

NUMSTEPS = 2
TIMESTEP_FRACTION = 1.0
NUM_STEPS = 20
BATCH_SIZE = 10
GUIDANCE_SCALE = 5.0
ETA = 1.0
NUM_PER_PROMPT = 7


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
        # if "base" not in source_model.name and source_model.ckpt is not None:
        #     pipe.unet.load_attn_procs(
        #         torch.load(io.BytesIO(source_model.ckpt.tobytes())), weights_only=True
        #     )

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
        sample_neg_prompt_embeds = neg_prompt_embed.repeat(BATCH_SIZE, 1, 1)

        # Sampling
        images, masks = zip(
            map(
                lambda base_asset: (base_asset.image, base_asset.mask),
                self.querier.get_random_base_assets_by_domain(
                    domain=source_model.domain, limit=BATCH_SIZE
                ),
            )
        )

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
                [[prompt] * BATCH_SIZE] * NUM_PER_PROMPT,
            )
        )

        sample_results = zip(
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
            post_latents.append(result[0])
            post_masks.append(result[0])
            post_mask_latents.append(result[0])

        latents = torch.stack(post_latents, dim=1)
        prompt_embeds = torch.stack(prompts_embeds, dim=1)
        mask_latents = torch.stack(post_mask_latents, dim=1)


def sample(pipe, prompt_embed, neg_prompt_embed, input_images, input_masks):
    images, _, latents, _, mask, masklatents = pipeline_with_logprob_inpaint(
        pipe,
        image=input_images,
        mask_image=input_masks,
        prompt_embeds=prompt_embed,
        negative_prompt_embeds=neg_prompt_embed,
        num_inference_steps=NUM_STEPS,
        guidance_scale=GUIDANCE_SCALE,
        eta=ETA,
        output_type="pt",
    )
    masklatents = torch.stack(masklatents, dim=1)
    images = images.cpu().detach()  # pyright: ignore
    latents = torch.stack(latents, dim=1).cpu().detach()  # pyright: ignore

    return images, latents, mask, masklatents
