from src.s3 import ImageUploader
from PIL import Image
from os import environ

image_uploader = ImageUploader(
    aws_access_key_id=environ["AWS_ACCESS_KEY_ID"],
    aws_secret_access_key=environ["AWS_SECRET_ACCESS_KEY"],
    bucket_name=environ["BUCKET_NAME"],
)

image = Image.open("./test/684B7A64-C283-4749-954D-8BF2D8A74192_1_105_c.jpeg")
url = "images/test2.jpg"
image_url = image_uploader.upload_image(image, url)
print(image_url)
