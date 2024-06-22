-- name: InsertSample :exec
INSERT INTO samples ( id, model_id )
VALUES ( $1, $2 );

-- name: GetModelByID :one
SELECT * FROM models
WHERE id = $1;

-- name: UpdateSampleFinished :exec
UPDATE samples
SET finished_at = now()
WHERE id = $1;

-- name: GetSampleByModelID :one
SELECT * FROM samples
WHERE model_id = $1
LIMIT 1;

-- name: InsertTrain :exec
INSERT INTO trains (id, sample_id)
VALUES ($1, $2);

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
    id, domain, parent_id, status,
    sample_id, train_id, updated_at
) VALUES (
    $1, $2, $3, $4,
    $5, $6, now()
);



-- name: ListAllUnfinishedSample :many
SELECT * FROM samples
WHERE finished_at IS NULL;

-- name: ListAllUnfinishedTrain :many
SELECT trains.id as train_id, m.* FROM trains
JOIN samples s ON sample_id = s.id
JOIN models m ON m.id = s.id
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





