import io

import numpy as np
import torch
from diffusers.loaders.utils import AttnProcsLayers
from diffusers.models.attention_processor import LoRAAttnProcessor
from diffusers.pipelines.stable_diffusion.pipeline_stable_diffusion_inpaint import (
    StableDiffusionInpaintPipeline,
)
from diffusers.schedulers.scheduling_ddim import DDIMScheduler
from PIL import Image

from src.database.query import Querier
from src.handlers.score import AestheticScorer


def prepare_pipe() -> StableDiffusionInpaintPipeline:
    pipe = StableDiffusionInpaintPipeline.from_pretrained(
        "runwayml/stable-diffusion-inpainting",
        torch_dtype=torch.float16,
    )
    pipe = pipe.to("cuda")
    pipe.vae.requires_grad_(False)
    pipe.text_encoder.requires_grad_(False)
    pipe.unet.requires_grad_(False)
    pipe.safety_checker = None
    pipe.set_progress_bar_config(
        position=1,
        leave=False,
        desc="Timestep",
        dynamic_ncols=True,
    )
    pipe.scheduler = DDIMScheduler.from_config(pipe.scheduler.config)
    inference_dtype = torch.float16
    pipe.vae.to("cuda", dtype=inference_dtype)
    pipe.text_encoder.to("cuda", dtype=inference_dtype)
    pipe.unet.to("cuda", dtype=inference_dtype)
    pipe = with_lora(pipeline=pipe)

    return pipe


def with_lora(
    pipeline: StableDiffusionInpaintPipeline,
) -> StableDiffusionInpaintPipeline:
    lora_attn_procs = {}
    for name in pipeline.unet.attn_processors.keys():
        cross_attention_dim = (
            None
            if name.endswith("attn1.processor")
            else pipeline.unet.config.cross_attention_dim
        )
        hidden_size = 0
        if name.startswith("mid_block"):
            hidden_size = pipeline.unet.config.block_out_channels[-1]
        elif name.startswith("up_blocks"):
            block_id = int(name[len("up_blocks.")])
            hidden_size = list(reversed(pipeline.unet.config.block_out_channels))[
                block_id
            ]
        elif name.startswith("down_blocks"):
            block_id = int(name[len("down_blocks.")])
            hidden_size = pipeline.unet.config.block_out_channels[block_id]

        lora_attn_procs[name] = LoRAAttnProcessor(
            hidden_size=hidden_size, cross_attention_dim=cross_attention_dim
        )

    pipeline.unet.set_attn_processor(lora_attn_procs)

    return pipeline


ADAM_BETA1 = 0.9
ADAM_BETA2 = 0.999
ADAM_WEIGHT_DECAY = 1e-4
ADAM_EPSILON = 1e-8
LEARNING_RATE = 3e-5


def prepare_optimizer(pipeline: StableDiffusionInpaintPipeline):
    trainable_layers = AttnProcsLayers(pipeline.unet.attn_processors)
    trainable_layers.backward()

    return torch.optim.AdamW(
        trainable_layers.parameters(),
        lr=LEARNING_RATE,
        betas=(ADAM_BETA1, ADAM_BETA2),
        weight_decay=ADAM_WEIGHT_DECAY,
        eps=ADAM_EPSILON,
    )


def light_reward():
    def _fn(images, _prompts, _metadata):
        reward = images.reshape(images.shape[0], -1).mean(1)
        return np.array(reward.cpu().detach()), {}

    return _fn


def jpeg_incompressibility():
    def _fn(images, _prompts, _metadata):
        if isinstance(images, torch.Tensor):
            images = (images * 255).round().clamp(0, 255).to(torch.uint8).cpu().numpy()
            images = images.transpose(0, 2, 3, 1)  # NCHW -> NHWC
        images = [Image.fromarray(image) for image in images]
        buffers = [io.BytesIO() for _ in images]
        for image, buffer in zip(images, buffers):
            image.save(buffer, format="JPEG", quality=95)
        sizes = [buffer.tell() / 1000 for buffer in buffers]
        return np.array(sizes), {}

    return _fn


def jpeg_compressibility():
    jpeg_fn = jpeg_incompressibility()

    def _fn(images, prompts, metadata):
        rew, meta = jpeg_fn(images, prompts, metadata)
        return -rew, meta

    return _fn


# need to setup table scorer before use this reward fn
def aesthetic_score(querier: Querier):
    scorer = AestheticScorer(dtype=torch.float32, querier=querier).cuda()

    def _fn(images, _):
        images = (images * 255).round().clamp(0, 255).to(torch.uint8)
        scores = scorer(images)
        return scores, {}

    return _fn


# return prompt and neg_prompt
def get_prompt(domain: str) -> tuple[str, str]:
    match domain:
        case "sessile":
            return "sessile", "pedunculated"
        case "pedunculated":
            return "pedunculated", "sessile"

    exit(1)


