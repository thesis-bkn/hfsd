-- name: InsertSample :exec
INSERT INTO samples ( id, model_id )
VALUES ( $1, $2 );

-- name: GetModelByID :one
SELECT * FROM models
WHERE id = $1;


-- name: UpdateModelStatus :exec
UPDATE models
SET status = $2
WHERE id = $1;

-- name: GetSampleByModelID :one
SELECT * FROM samples
WHERE model_id = $1
LIMIT 1;

-- name: CheckSampleFinishedByModelID :one
SELECT EXISTS (
    SELECT *
    FROM samples
    WHERE model_id = $1
    AND finished_at IS NOT NULL
);

-- name: ListModelByDomain :many
SELECT *
FROM models
WHERE domain = $1;

-- name: InsertTrain :exec
INSERT INTO trains (id, model_id, sample_id)
VALUES ($1, $2, $3);

-- name: InsertInference :exec
INSERT INTO inferences (
    id, model_id, prompt, neg_prompt
) VALUES ( $1, $2, $3, $4 );

-- name: CheckModelExists :one
SELECT EXISTS (
    SELECT *
    FROM models
    WHERE id = $1
);

-- name: InsertModel :exec
INSERT INTO models (
    id, domain, name, parent_id, status,
    sample_id, train_id, updated_at
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, now()
);



-- name: ListAllUnfinishedSample :many
SELECT * FROM samples
WHERE finished_at IS NULL;

-- name: ListAllUnfinishedTrain :many
SELECT trains.id as train_id, m.* FROM trains
JOIN models m ON m.id = trains.id
WHERE trains.finished_at IS NULL;

-- name: ListAllUnfinishedInferences :many
SELECT  i.id as inference_id,
        i.prompt,
        i.neg_prompt,
        i.finished_at,
        m.*
FROM inferences i
JOIN models m on m.id = i.model_id
WHERE finished_at IS NULL;

-- name: ListModels :many
SELECT * FROM models
WHERE id = ANY($1::text[]);

-- name: ListFinishedInferences :many
SELECT * FROM inferences i
JOIN models m ON i.model_id = m.id
WHERE finished_at IS NOT NULL
LIMIT $1 OFFSET $2;

-- name: UpdateInferenceFinished :exec
UPDATE inferences
SET finished_at = now()
WHERE id = $1;

-- name: UpdateSampleFinished :exec
UPDATE samples
SET finished_at = now()
WHERE id = $1;

-- name: UpdateTrainFinished :exec
UPDATE trains
SET finished_at = now()
WHERE id = $1;

-- name: InsertTask :one
INSERT INTO tasks (
    task_type, content,
    status, estimate, updated_at
)
VALUES ( $1, $2, $3, $4, now())
RETURNING task_id;

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET status = $2,
    estimate = $3,
    updated_at = now()
WHERE task_id = $1;


-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY created_at DESC;
