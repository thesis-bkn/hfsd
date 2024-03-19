from src.database.query import InsertBaseAssetParams
from src.s3 import ImageUploader
from PIL import Image
from os import environ
from sqlalchemy.engine import create_engine
from src.config import ConfigReader
from src.database import Querier
from tqdm import tqdm

import os

image_uploader = ImageUploader(
    aws_access_key_id=environ["AWS_ACCESS_KEY_ID"],
    aws_secret_access_key=environ["AWS_SECRET_ACCESS_KEY"],
    bucket_name=environ["BUCKET_NAME"],
    endpoint_url=environ["S3_ENDPOINT_URL"],
)
config = ConfigReader("config.yml")
conn = create_engine(
    environ["DATABASE_URL"].replace("postgres://", "postgresql://")
).connect()
querier = Querier(conn)

base_assets = config.get_value("base_assets")
for domain in base_assets:
    for filename in tqdm(os.listdir(domain["image_dir"])):
        id, _ = os.path.splitext(filename)
        image = Image.open(os.path.join(domain["image_dir"], filename))
        mask = Image.open(os.path.join(domain["mask_dir"], filename))

        params = InsertBaseAssetParams(
            id=id,
            image=memoryview(image.tobytes()),
            image_url=os.path.join("images", filename),
            mask=memoryview(mask.tobytes()),
            mask_url=os.path.join("masks", filename),
            domain=domain["domain"],
        )
        image_uploader.upload_image(image, params.image_url)
        image_uploader.upload_image(mask, params.mask_url)
        querier.insert_base_asset(params)

    conn.commit()
