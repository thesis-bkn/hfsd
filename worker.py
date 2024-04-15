import datetime
import time
from os import environ
from typing import Optional

import schedule
from sqlalchemy.engine import create_engine

from src.config import ConfigReader
from src.database import Querier
from src.database.models import Task, TaskVariant
from src.handlers import InferenceHandler
from src.handlers.finetune import FinetuneHandler
from src.handlers.sample import SampleHander
from src.s3 import ImageUploader


class CronJob:
    def __init__(self):
        # Initialize the schedule
        self.schedule = schedule
        self.schedule.every().second.do(self.run_job)

        self.config = ConfigReader("config.yml")
        self.conn = create_engine(
            environ["DATABASE_URL"].replace("postgres://", "postgresql://")
        ).connect()
        self.querier = Querier(self.conn)
        self.current_task: Optional[Task] = None
        self.uploader = ImageUploader(
            aws_access_key_id=environ["AWS_ACCESS_KEY_ID"],
            aws_secret_access_key=environ["AWS_SECRET_ACCESS_KEY"],
            bucket_name=environ["BUCKET_NAME"],
            endpoint_url=environ["S3_ENDPOINT_URL"],
        )

        self.inf_handler = InferenceHandler(conn=self.conn, uploader=self.uploader)
        self.sample_handler = SampleHander(conn=self.conn, uploader=self.uploader)
        self.finetune_handler = FinetuneHandler(conn=self.conn, uploader=self.uploader)

    def run_job(self):
        # Define the job to be executed
        print("Cron job is running at:", time.strftime("%Y-%m-%d %H:%M:%S"))

        # Check is there any task need to be done
        pending_task = self.querier.get_earliest_pending_task()
        if pending_task is None:
            return

        # mark that this task is being handled
        self.querier.update_task_status(
            id=pending_task.id,
            handled_at=datetime.datetime.now(datetime.UTC),
            finished_at=None,
        )
        self.conn.commit()

        # Specify what task is this
        match pending_task.task_type:
            case TaskVariant.INFERENCE:
                self.inf_handler.handle(task=pending_task)
            case TaskVariant.SAMPLE:
                self.sample_handler.handle(task=pending_task)
            case TaskVariant.FINETUNE:
                self.finetune_handler.handle(task=pending_task)


    def start(self):
        # Start the scheduler
        while True:
            self.schedule.run_pending()
            time.sleep(1)


if __name__ == "__main__":
    cron = CronJob()
    cron.start()
