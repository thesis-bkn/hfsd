-- name: GetModel :one
SELECT * FROM models
WHERE id = $1 LIMIT 1;

-- name: ListModelsByDomain :many
SELECT * FROM models
WHERE domain = $1;

-- name: InsertModel :exec
INSERT INTO models (id, domain, name, base, ckpt)
VALUES ($1, $2, $3, $4, $5);

-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1 AND task_type = $2
LIMIT 1;
