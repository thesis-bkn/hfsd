from typing import Optional
from src import database
from src.config import ConfigReader
from src.database import Querier
from sqlalchemy.engine import create_engine
from os import environ
import schedule
import time
import datetime

from src.database.models import Task, TaskVariant


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

    def run_job(self):
        # Define the job to be executed
        print("Cron job is running at:", time.strftime("%Y-%m-%d %H:%M:%S"))

        # Check is there any task need to be done
        pending_task = self.querier.get_earliest_pending_task()
        if pending_task is None:
            return

        # Specify what task is this
        match pending_task.task_type:
            case TaskVariant.INFERENCE:
                print("inference task")
            case TaskVariant.FINETUNE:
                print("fine tune task")
            case TaskVariant.SAMPLE:
                print("sample task")

        # mark that this task is being handled
        self.querier.update_task_status(
            id=pending_task.id,
            task_type=pending_task.task_type,
            handled_at=datetime.datetime.now(datetime.UTC),
            finished_at=None,
        )

    def start(self):
        # Start the scheduler
        while True:
            self.schedule.run_pending()
            time.sleep(1)


if __name__ == "__main__":
    cron = CronJob()
    cron.start()
