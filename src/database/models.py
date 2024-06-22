# Code generated by sqlc. DO NOT EDIT.
# versions:
#   sqlc v1.26.0
import dataclasses
import datetime
from typing import Optional


@dataclasses.dataclass()
class Inference:
    id: str
    model_id: str
    prompt: str
    neg_prompt: str
    finished_at: Optional[datetime.datetime]


@dataclasses.dataclass()
class Model:
    id: str
    domain: str
    parent_id: Optional[str]
    status: str
    sample_id: Optional[str]
    train_id: Optional[str]
    updated_at: Optional[datetime.datetime]
    created_at: Optional[datetime.datetime]


@dataclasses.dataclass()
class Sample:
    id: str
    model_id: str
    finished_at: Optional[datetime.datetime]
    created_at: datetime.datetime


@dataclasses.dataclass()
class Train:
    id: str
    sample_id: str
    created_at: datetime.datetime
    finished_at: Optional[datetime.datetime]
