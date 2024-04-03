# Code generated by sqlc. DO NOT EDIT.
# versions:
#   sqlc v1.25.0
import dataclasses
import datetime
import enum
from typing import Any, Optional


class TaskVariant(str, enum.Enum):
    INFERENCE = "inference"
    SAMPLE = "sample"
    FINETUNE = "finetune"


@dataclasses.dataclass()
class Asset:
    task_id: str
    order: int
    prompt: Optional[str]
    image: memoryview
    image_url: str
    mask: memoryview
    mask_url: str


@dataclasses.dataclass()
class BaseAsset:
    id: str
    image: memoryview
    image_url: str
    mask: memoryview
    mask_url: str
    domain: str


@dataclasses.dataclass()
class Inference:
    id: str
    prompt: Optional[str]
    image: memoryview
    image_url: str
    mask: memoryview
    mask_url: str
    output: memoryview
    output_url: str
    from_model: str


@dataclasses.dataclass()
class Model:
    id: str
    domain: str
    name: str
    base: str
    ckpt: memoryview
    created_at: Optional[datetime.datetime]


@dataclasses.dataclass()
class Task:
    id: str
    source_model_id: str
    output_model_id: Optional[str]
    task_type: TaskVariant
    created_at: Optional[datetime.datetime]
    handled_at: Optional[datetime.datetime]
    finished_at: Optional[datetime.datetime]
    human_prefs: Optional[Any]
    prompt_embeds: Optional[memoryview]
    latents: Optional[memoryview]
    timesteps: Optional[memoryview]
    next_latents: Optional[memoryview]
    image_torchs: Optional[memoryview]
