import datetime
import io
import os

import torch
from diffusers.pipelines.stable_diffusion.pipeline_stable_diffusion_inpaint import (
    StableDiffusionInpaintPipeline,
)
from PIL import Image
from sqlalchemy.engine.base import Connection

from src.database.models import Task
from src.database.query import Querier, SaveInferenceParams
from src.s3 import ImageUploader


class InferenceHandler:
    def __init__(self, conn: Connection, uploader: ImageUploader) -> None:
        self.conn = conn
        self.uploader = uploader
        self.querier = Querier(self.conn)

    def handle(self, task: Task):
        pipe = StableDiffusionInpaintPipeline.from_pretrained(
            "runwayml/stable-diffusion-inpainting",
            torch_dtype=torch.float16,
        )
        pipe = pipe.to("cuda")
        pipe.vae.requires_grad_(False)
        pipe.text_encoder.requires_grad_(False)
        pipe.unet.requires_grad_(False)

        if "base" not in task.source_model_id:
            model = self.querier.get_model(id=task.source_model_id)
            if model is None:
                exit("cannot find model OMG")

            pipe.unet.load_attn_procs(
                torch.load(io.BytesIO(model.ckpt.tobytes())), weights_only=True
            )

        # get image
        asset = self.querier.get_first_asset_by_model_id(task_id=task.id)
        if asset is None:
            raise Exception("can not find asset to inference")
        image = Image.open(io.BytesIO(asset.image.tobytes()))
        mask = Image.open(io.BytesIO(asset.mask.tobytes()))  # pyright: ignore
        output = pipe(prompt=asset.prompt, image=image, mask_image=mask).images[0]

        # upload output image to cloudfly
        # then save new record inside table inference
        # update status of task
        image_bytes = io.BytesIO()
        output.save(image_bytes, format="JPEG")
        image_bytes.seek(0)  # Rewind to the beginning of the BytesIO buffer
        output_bytes = image_bytes.getvalue()

        key = os.path.join("output", str(task.id))
        self.uploader.upload_image(
            output_bytes,
            s3_key=key,
        )

        self.querier.save_inference(
            SaveInferenceParams(
                id="inference" + str(task.id),
                prompt=asset.prompt,
                image=memoryview(asset.image),
                image_url=asset.image_url,
                mask=memoryview(asset.mask),  # pyright: ignore
                mask_url=asset.mask_url,
                output=memoryview(output_bytes),
                output_url=key,
                from_model=task.source_model_id,
            )
        )

        self.querier.update_task_status(
            id=task.id,
            handled_at=None,
            finished_at=datetime.datetime.now(datetime.UTC),
        )

        self.conn.commit()
