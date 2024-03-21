import boto3
from typing import Optional
import io
import os
from PIL import Image


class ImageUploader:
    _instance = None

    def __new__(cls, *_, **__):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance.__initialized = False
        return cls._instance

    def __init__(
        self,
        aws_access_key_id: str,
        aws_secret_access_key: str,
        bucket_name: str,
        endpoint_url: str,
    ):
        if self.__initialized:
            return
        self.__initialized = True
        self.s3_client = boto3.client(
            "s3",
            endpoint_url=endpoint_url,
            aws_access_key_id=aws_access_key_id,
            aws_secret_access_key=aws_secret_access_key,
        )
        self.bucket_name = bucket_name

    def upload_image(self, image: bytes, s3_key: str) -> Optional[str]:
        try:
            # Convert PIL image to bytes
            image_bytes = io.BytesIO(image)

            # Upload image to S3
            self.s3_client.upload_fileobj(image_bytes, self.bucket_name, s3_key)

            # Generate pre-signed URL for the uploaded image
            url = self.generate_presigned_url(s3_key)
            return url
        except Exception as e:
            return None

    def generate_presigned_url(self, s3_key, expiration=3600) -> Optional[str]:
        try:
            url = self.s3_client.generate_presigned_url(
                "get_object",
                Params={"Bucket": self.bucket_name, "Key": s3_key},
                ExpiresIn=expiration,
            )
            return url
        except Exception as e:
            print(f"Error generating presigned URL: {e}")
            return None
