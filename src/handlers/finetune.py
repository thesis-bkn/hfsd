from sqlalchemy.engine.base import Connection
from src.database.query import Querier
from src.s3 import ImageUploader


class FinetuneHandler:
    def __init__(self, conn: Connection, uploader: ImageUploader) -> None:
        self.conn = conn
        self.uploader = uploader
        self.querier = Querier(self.conn)
