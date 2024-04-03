// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getEarliestPendingTask = `-- name: GetEarliestPendingTask :one
SELECT id, source_model_id, output_model_id, task_type, created_at, handled_at, finished_at, human_prefs, prompt_embeds, latents, timesteps, next_latents, image_torchs FROM tasks
WHERE handled_at IS NULL
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetEarliestPendingTask(ctx context.Context) (Task, error) {
	row := q.db.QueryRow(ctx, getEarliestPendingTask)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.SourceModelID,
		&i.OutputModelID,
		&i.TaskType,
		&i.CreatedAt,
		&i.HandledAt,
		&i.FinishedAt,
		&i.HumanPrefs,
		&i.PromptEmbeds,
		&i.Latents,
		&i.Timesteps,
		&i.NextLatents,
		&i.ImageTorchs,
	)
	return i, err
}

const getFirstAssetByModelID = `-- name: GetFirstAssetByModelID :one
SELECT task_id, "order", prompt, image, image_url, mask, mask_url FROM assets
WHERE task_id = $1
AND "order" = 0
LIMIT 1
`

func (q *Queries) GetFirstAssetByModelID(ctx context.Context, taskID string) (Asset, error) {
	row := q.db.QueryRow(ctx, getFirstAssetByModelID, taskID)
	var i Asset
	err := row.Scan(
		&i.TaskID,
		&i.Order,
		&i.Prompt,
		&i.Image,
		&i.ImageUrl,
		&i.Mask,
		&i.MaskUrl,
	)
	return i, err
}

const getModel = `-- name: GetModel :one
SELECT id, domain, name, base, ckpt, created_at FROM models
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetModel(ctx context.Context, id string) (Model, error) {
	row := q.db.QueryRow(ctx, getModel, id)
	var i Model
	err := row.Scan(
		&i.ID,
		&i.Domain,
		&i.Name,
		&i.Base,
		&i.Ckpt,
		&i.CreatedAt,
	)
	return i, err
}

const getTask = `-- name: GetTask :one
SELECT id, source_model_id, output_model_id, task_type, created_at, handled_at, finished_at, human_prefs, prompt_embeds, latents, timesteps, next_latents, image_torchs FROM tasks
WHERE id = $1 AND task_type = $2
LIMIT 1
`

type GetTaskParams struct {
	ID       string
	TaskType TaskVariant
}

func (q *Queries) GetTask(ctx context.Context, arg GetTaskParams) (Task, error) {
	row := q.db.QueryRow(ctx, getTask, arg.ID, arg.TaskType)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.SourceModelID,
		&i.OutputModelID,
		&i.TaskType,
		&i.CreatedAt,
		&i.HandledAt,
		&i.FinishedAt,
		&i.HumanPrefs,
		&i.PromptEmbeds,
		&i.Latents,
		&i.Timesteps,
		&i.NextLatents,
		&i.ImageTorchs,
	)
	return i, err
}

const insertAsset = `-- name: InsertAsset :exec
INSERT INTO assets (task_id, "order", prompt, image, image_url, mask, mask_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type InsertAssetParams struct {
	TaskID   string
	Order    int16
	Prompt   string
	Image    []byte
	ImageUrl string
	Mask     []byte
	MaskUrl  string
}

func (q *Queries) InsertAsset(ctx context.Context, arg InsertAssetParams) error {
	_, err := q.db.Exec(ctx, insertAsset,
		arg.TaskID,
		arg.Order,
		arg.Prompt,
		arg.Image,
		arg.ImageUrl,
		arg.Mask,
		arg.MaskUrl,
	)
	return err
}

const insertBaseAsset = `-- name: InsertBaseAsset :exec
INSERT INTO base_assets (id, image, image_url, mask, mask_url, domain)
VALUES ($1, $2, $3, $4, $5, $6)
`

type InsertBaseAssetParams struct {
	ID       string
	Image    []byte
	ImageUrl string
	Mask     []byte
	MaskUrl  string
	Domain   string
}

func (q *Queries) InsertBaseAsset(ctx context.Context, arg InsertBaseAssetParams) error {
	_, err := q.db.Exec(ctx, insertBaseAsset,
		arg.ID,
		arg.Image,
		arg.ImageUrl,
		arg.Mask,
		arg.MaskUrl,
		arg.Domain,
	)
	return err
}

const insertInferenceTask = `-- name: InsertInferenceTask :exec
INSERT INTO tasks (id, source_model_id, task_type)
VALUES ( $1, $2, 'inference' )
`

type InsertInferenceTaskParams struct {
	ID            string
	SourceModelID string
}

func (q *Queries) InsertInferenceTask(ctx context.Context, arg InsertInferenceTaskParams) error {
	_, err := q.db.Exec(ctx, insertInferenceTask, arg.ID, arg.SourceModelID)
	return err
}

const insertModel = `-- name: InsertModel :exec
INSERT INTO models (id, domain, name, base, ckpt)
VALUES ($1, $2, $3, $4, $5)
`

type InsertModelParams struct {
	ID     string
	Domain string
	Name   string
	Base   string
	Ckpt   []byte
}

func (q *Queries) InsertModel(ctx context.Context, arg InsertModelParams) error {
	_, err := q.db.Exec(ctx, insertModel,
		arg.ID,
		arg.Domain,
		arg.Name,
		arg.Base,
		arg.Ckpt,
	)
	return err
}

const listAllTaskWithAsset = `-- name: ListAllTaskWithAsset :many
SELECT tasks.id, tasks.source_model_id, tasks.output_model_id, tasks.task_type, tasks.created_at, tasks.handled_at, tasks.finished_at, tasks.human_prefs, tasks.prompt_embeds, tasks.latents, tasks.timesteps, tasks.next_latents, tasks.image_torchs, assets.task_id, assets."order", assets.prompt, assets.image, assets.image_url, assets.mask, assets.mask_url
FROM tasks
JOIN assets ON assets.task_id = tasks.id
WHERE assets."order" = 0
ORDER BY tasks.created_at DESC
`

type ListAllTaskWithAssetRow struct {
	Task  Task
	Asset Asset
}

func (q *Queries) ListAllTaskWithAsset(ctx context.Context) ([]ListAllTaskWithAssetRow, error) {
	rows, err := q.db.Query(ctx, listAllTaskWithAsset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAllTaskWithAssetRow
	for rows.Next() {
		var i ListAllTaskWithAssetRow
		if err := rows.Scan(
			&i.Task.ID,
			&i.Task.SourceModelID,
			&i.Task.OutputModelID,
			&i.Task.TaskType,
			&i.Task.CreatedAt,
			&i.Task.HandledAt,
			&i.Task.FinishedAt,
			&i.Task.HumanPrefs,
			&i.Task.PromptEmbeds,
			&i.Task.Latents,
			&i.Task.Timesteps,
			&i.Task.NextLatents,
			&i.Task.ImageTorchs,
			&i.Asset.TaskID,
			&i.Asset.Order,
			&i.Asset.Prompt,
			&i.Asset.Image,
			&i.Asset.ImageUrl,
			&i.Asset.Mask,
			&i.Asset.MaskUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listModelsByDomain = `-- name: ListModelsByDomain :many
SELECT id, domain, name, base, ckpt, created_at FROM models
WHERE domain = $1
`

func (q *Queries) ListModelsByDomain(ctx context.Context, domain string) ([]Model, error) {
	rows, err := q.db.Query(ctx, listModelsByDomain, domain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Model
	for rows.Next() {
		var i Model
		if err := rows.Scan(
			&i.ID,
			&i.Domain,
			&i.Name,
			&i.Base,
			&i.Ckpt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveInference = `-- name: SaveInference :exec
INSERT INTO inferences (id, prompt, image, image_url, mask, mask_url, output, output_url, from_model)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

type SaveInferenceParams struct {
	ID        string
	Prompt    pgtype.Text
	Image     []byte
	ImageUrl  string
	Mask      []byte
	MaskUrl   string
	Output    []byte
	OutputUrl string
	FromModel string
}

func (q *Queries) SaveInference(ctx context.Context, arg SaveInferenceParams) error {
	_, err := q.db.Exec(ctx, saveInference,
		arg.ID,
		arg.Prompt,
		arg.Image,
		arg.ImageUrl,
		arg.Mask,
		arg.MaskUrl,
		arg.Output,
		arg.OutputUrl,
		arg.FromModel,
	)
	return err
}

const updateTaskStatus = `-- name: UpdateTaskStatus :exec
INSERT INTO tasks (id, source_model_id, task_type, handled_at, finished_at)
VALUES ($1, '', $2, $3, $4)
ON CONFLICT (id, task_type)
DO UPDATE SET
    handled_at = COALESCE(tasks.handled_at, EXCLUDED.handled_at),
    finished_at = COALESCE(tasks.finished_at, EXCLUDED.finished_at)
`

type UpdateTaskStatusParams struct {
	ID         string
	TaskType   TaskVariant
	HandledAt  pgtype.Timestamp
	FinishedAt pgtype.Timestamp
}

func (q *Queries) UpdateTaskStatus(ctx context.Context, arg UpdateTaskStatusParams) error {
	_, err := q.db.Exec(ctx, updateTaskStatus,
		arg.ID,
		arg.TaskType,
		arg.HandledAt,
		arg.FinishedAt,
	)
	return err
}
