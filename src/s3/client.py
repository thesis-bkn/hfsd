import boto3
from typing import Optional
import io
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
        endpoint_url: str = "https://s3.cloudfly.vn",
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

    def upload_image(self, pil_image: Image.Image, s3_key: str) -> Optional[str]:
        try:
            # Convert PIL image to bytes
            image_bytes = io.BytesIO()
            pil_image.save(image_bytes, format="JPEG")
            image_bytes.seek(0)

            # Upload image to S3
            self.s3_client.upload_fileobj(image_bytes, self.bucket_name, s3_key)

            # Generate pre-signed URL for the uploaded image
            url = self.generate_presigned_url(s3_key)
            print(f"Image uploaded successfully to S3 with key: {s3_key}")
            return url
        except Exception as e:
            print(f"Error uploading image to S3: {e}")
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


# # Usage example:
# # Initialize the uploader
# uploader = S3ImageUploader(
#     aws_access_key_id="YOUR_AWS_ACCESS_KEY_ID",
#     aws_secret_access_key="YOUR_AWS_SECRET_ACCESS_KEY",
#     bucket_name="YOUR_S3_BUCKET_NAME",
# )
#
# # Load PIL image from file or any source
# pil_image = Image.open("path/to/local/image.jpg")
#
# # Specify S3 key
# s3_key = "images/image.jpg"
#
# # Upload the image to S3 and get the presigned URL
# image_url = uploader.upload_image(pil_image, s3_key)
#
# if image_url:
#     print("Presigned URL:", image_url)
