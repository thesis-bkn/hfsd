-- name: GetModel :one
SELECT * FROM models
WHERE id = $1 LIMIT 1;

-- name: ListModelsByDomain :many
SELECT * FROM models
WHERE domain = $1;

-- name: GetEarliestPendingTask :one
SELECT * FROM tasks
WHERE handled_at IS NULL
ORDER BY created_at DESC
LIMIT 1;

-- name: ListAllTaskWithAsset :many
SELECT sqlc.embed(tasks), sqlc.embed(assets)
FROM tasks
JOIN assets ON assets.task_id = tasks.id
WHERE assets."order" = 0
ORDER BY tasks.created_at DESC;

-- name: InsertModel :exec
INSERT INTO models (id, domain, name, base, ckpt)
VALUES ($1, $2, $3, $4, $5);

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1 AND task_type = $2
LIMIT 1;

-- name: InsertBaseAsset :exec
INSERT INTO base_assets (id, image, image_url, mask, mask_url, domain)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: InsertInferenceTask :exec
INSERT INTO tasks (id, source_model_id, task_type)
VALUES ( $1, $2, 'inference' );

-- name: InsertAsset :exec
INSERT INTO assets (task_id, "order", prompt, image, image_url, mask, mask_url)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateTaskStatus :exec
INSERT INTO tasks (id, task_type, handled_at, finished_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id, task_type)
DO UPDATE SET
    handled_at = COALESCE(tasks.handled_at, EXCLUDED.handled_at),
    finished_at = COALESCE(tasks.finished_at, EXCLUDED.finished_at);

-- name: GetFirstAssetByModelID :one
SELECT * FROM assets
WHERE task_id = $1
AND "order" = 0
LIMIT 1;

-- name: SaveInference :exec
INSERT INTO inferences (id, prompt, image, image_url, mask, mask_url, output, output_url, from_model)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

