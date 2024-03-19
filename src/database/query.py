# Code generated by sqlc. DO NOT EDIT.
# versions:
#   sqlc v1.25.0
# source: query.sql
import dataclasses
from typing import AsyncIterator, Iterator, Optional

import sqlalchemy
import sqlalchemy.ext.asyncio

from database import models


GET_MODEL = """-- name: get_model \\:one
SELECT id, domain, name, base, ckpt, created_at FROM models
WHERE id = :p1 LIMIT 1
"""


GET_TASK = """-- name: get_task \\:one
SELECT id, source_model_id, output_model_id, task_type, created_at, handled_at, finished_at, human_prefs, prompt_embeds, latents, timesteps, next_latents, image_torchs FROM tasks
WHERE id = :p1 AND task_type = :p2
LIMIT 1
"""


INSERT_MODEL = """-- name: insert_model \\:exec
INSERT INTO models (id, domain, name, base, ckpt)
VALUES (:p1, :p2, :p3, :p4, :p5)
"""


@dataclasses.dataclass()
class InsertModelParams:
    id: str
    domain: str
    name: str
    base: str
    ckpt: memoryview


LIST_MODELS_BY_DOMAIN = """-- name: list_models_by_domain \\:many
SELECT id, domain, name, base, ckpt, created_at FROM models
WHERE domain = :p1
"""


class Querier:
    def __init__(self, conn: sqlalchemy.engine.Connection):
        self._conn = conn

    def get_model(self, *, id: str) -> Optional[models.Model]:
        row = self._conn.execute(sqlalchemy.text(GET_MODEL), {"p1": id}).first()
        if row is None:
            return None
        return models.Model(
            id=row[0],
            domain=row[1],
            name=row[2],
            base=row[3],
            ckpt=row[4],
            created_at=row[5],
        )

    def get_task(self, *, id: str, task_type: str) -> Optional[models.Task]:
        row = self._conn.execute(sqlalchemy.text(GET_TASK), {"p1": id, "p2": task_type}).first()
        if row is None:
            return None
        return models.Task(
            id=row[0],
            source_model_id=row[1],
            output_model_id=row[2],
            task_type=row[3],
            created_at=row[4],
            handled_at=row[5],
            finished_at=row[6],
            human_prefs=row[7],
            prompt_embeds=row[8],
            latents=row[9],
            timesteps=row[10],
            next_latents=row[11],
            image_torchs=row[12],
        )

    def insert_model(self, arg: InsertModelParams) -> None:
        self._conn.execute(sqlalchemy.text(INSERT_MODEL), {
            "p1": arg.id,
            "p2": arg.domain,
            "p3": arg.name,
            "p4": arg.base,
            "p5": arg.ckpt,
        })

    def list_models_by_domain(self, *, domain: str) -> Iterator[models.Model]:
        result = self._conn.execute(sqlalchemy.text(LIST_MODELS_BY_DOMAIN), {"p1": domain})
        for row in result:
            yield models.Model(
                id=row[0],
                domain=row[1],
                name=row[2],
                base=row[3],
                ckpt=row[4],
                created_at=row[5],
            )


class AsyncQuerier:
    def __init__(self, conn: sqlalchemy.ext.asyncio.AsyncConnection):
        self._conn = conn

    async def get_model(self, *, id: str) -> Optional[models.Model]:
        row = (await self._conn.execute(sqlalchemy.text(GET_MODEL), {"p1": id})).first()
        if row is None:
            return None
        return models.Model(
            id=row[0],
            domain=row[1],
            name=row[2],
            base=row[3],
            ckpt=row[4],
            created_at=row[5],
        )

    async def get_task(self, *, id: str, task_type: str) -> Optional[models.Task]:
        row = (await self._conn.execute(sqlalchemy.text(GET_TASK), {"p1": id, "p2": task_type})).first()
        if row is None:
            return None
        return models.Task(
            id=row[0],
            source_model_id=row[1],
            output_model_id=row[2],
            task_type=row[3],
            created_at=row[4],
            handled_at=row[5],
            finished_at=row[6],
            human_prefs=row[7],
            prompt_embeds=row[8],
            latents=row[9],
            timesteps=row[10],
            next_latents=row[11],
            image_torchs=row[12],
        )

    async def insert_model(self, arg: InsertModelParams) -> None:
        await self._conn.execute(sqlalchemy.text(INSERT_MODEL), {
            "p1": arg.id,
            "p2": arg.domain,
            "p3": arg.name,
            "p4": arg.base,
            "p5": arg.ckpt,
        })

    async def list_models_by_domain(self, *, domain: str) -> AsyncIterator[models.Model]:
        result = await self._conn.stream(sqlalchemy.text(LIST_MODELS_BY_DOMAIN), {"p1": domain})
        async for row in result:
            yield models.Model(
                id=row[0],
                domain=row[1],
                name=row[2],
                base=row[3],
                ckpt=row[4],
                created_at=row[5],
            )
