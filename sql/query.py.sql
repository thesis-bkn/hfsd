-- name: GetModel :one
SELECT * FROM models
WHERE id = $1 LIMIT 1;
