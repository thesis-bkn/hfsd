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

-- name: ListAllTask :many
SELECT * FROM tasks;

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
INSERT INTO assets (task_id, "order", image, image_url, mask, mask_url)
VALUES ($1, $2, $3, $4, $5, $6);
