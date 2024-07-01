// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkModelExists = `-- name: CheckModelExists :one
SELECT EXISTS (
    SELECT id, domain, name, parent_id, status, sample_id, train_id, updated_at, created_at
    FROM models
    WHERE id = $1
)
`

func (q *Queries) CheckModelExists(ctx context.Context, id string) (bool, error) {
	row := q.db.QueryRow(ctx, checkModelExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkSampleFinishedByModelID = `-- name: CheckSampleFinishedByModelID :one
SELECT EXISTS (
    SELECT id, model_id, finished_at, created_at
    FROM samples
    WHERE model_id = $1
    AND finished_at IS NOT NULL
)
`

func (q *Queries) CheckSampleFinishedByModelID(ctx context.Context, modelID string) (bool, error) {
	row := q.db.QueryRow(ctx, checkSampleFinishedByModelID, modelID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getModelByID = `-- name: GetModelByID :one
SELECT id, domain, name, parent_id, status, sample_id, train_id, updated_at, created_at FROM models
WHERE id = $1
`

func (q *Queries) GetModelByID(ctx context.Context, id string) (Model, error) {
	row := q.db.QueryRow(ctx, getModelByID, id)
	var i Model
	err := row.Scan(
		&i.ID,
		&i.Domain,
		&i.Name,
		&i.ParentID,
		&i.Status,
		&i.SampleID,
		&i.TrainID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSampleByModelID = `-- name: GetSampleByModelID :one
SELECT id, model_id, finished_at, created_at FROM samples
WHERE model_id = $1
LIMIT 1
`

func (q *Queries) GetSampleByModelID(ctx context.Context, modelID string) (Sample, error) {
	row := q.db.QueryRow(ctx, getSampleByModelID, modelID)
	var i Sample
	err := row.Scan(
		&i.ID,
		&i.ModelID,
		&i.FinishedAt,
		&i.CreatedAt,
	)
	return i, err
}

const insertInference = `-- name: InsertInference :exec
INSERT INTO inferences (
    id, model_id, prompt, neg_prompt
) VALUES ( $1, $2, $3, $4 )
`

type InsertInferenceParams struct {
	ID        string
	ModelID   string
	Prompt    string
	NegPrompt string
}

func (q *Queries) InsertInference(ctx context.Context, arg InsertInferenceParams) error {
	_, err := q.db.Exec(ctx, insertInference,
		arg.ID,
		arg.ModelID,
		arg.Prompt,
		arg.NegPrompt,
	)
	return err
}

const insertModel = `-- name: InsertModel :exec
INSERT INTO models (
    id, domain, name, parent_id, status,
    sample_id, train_id, updated_at
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, now()
)
`

type InsertModelParams struct {
	ID       string
	Domain   string
	Name     string
	ParentID pgtype.Text
	Status   string
	SampleID pgtype.Text
	TrainID  pgtype.Text
}

func (q *Queries) InsertModel(ctx context.Context, arg InsertModelParams) error {
	_, err := q.db.Exec(ctx, insertModel,
		arg.ID,
		arg.Domain,
		arg.Name,
		arg.ParentID,
		arg.Status,
		arg.SampleID,
		arg.TrainID,
	)
	return err
}

const insertSample = `-- name: InsertSample :exec
INSERT INTO samples ( id, model_id )
VALUES ( $1, $2 )
`

type InsertSampleParams struct {
	ID      string
	ModelID string
}

func (q *Queries) InsertSample(ctx context.Context, arg InsertSampleParams) error {
	_, err := q.db.Exec(ctx, insertSample, arg.ID, arg.ModelID)
	return err
}

const insertTask = `-- name: InsertTask :one
INSERT INTO tasks (
    task_type, content,
    status, estimate, updated_at
)
VALUES ( $1, $2, $3, $4, now())
RETURNING task_id
`

type InsertTaskParams struct {
	TaskType string
	Content  string
	Status   string
	Estimate int64
}

func (q *Queries) InsertTask(ctx context.Context, arg InsertTaskParams) (int32, error) {
	row := q.db.QueryRow(ctx, insertTask,
		arg.TaskType,
		arg.Content,
		arg.Status,
		arg.Estimate,
	)
	var task_id int32
	err := row.Scan(&task_id)
	return task_id, err
}

const insertTrain = `-- name: InsertTrain :exec
INSERT INTO trains (id, model_id, sample_id)
VALUES ($1, $2, $3)
`

type InsertTrainParams struct {
	ID       string
	ModelID  string
	SampleID string
}

func (q *Queries) InsertTrain(ctx context.Context, arg InsertTrainParams) error {
	_, err := q.db.Exec(ctx, insertTrain, arg.ID, arg.ModelID, arg.SampleID)
	return err
}

const listAllUnfinishedInferences = `-- name: ListAllUnfinishedInferences :many
SELECT  i.id as inference_id,
        i.prompt,
        i.neg_prompt,
        i.finished_at,
        m.id, m.domain, m.name, m.parent_id, m.status, m.sample_id, m.train_id, m.updated_at, m.created_at
FROM inferences i
JOIN models m on m.id = i.model_id
WHERE finished_at IS NULL
`

type ListAllUnfinishedInferencesRow struct {
	InferenceID string
	Prompt      string
	NegPrompt   string
	FinishedAt  pgtype.Timestamptz
	ID          string
	Domain      string
	Name        string
	ParentID    pgtype.Text
	Status      string
	SampleID    pgtype.Text
	TrainID     pgtype.Text
	UpdatedAt   pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
}

func (q *Queries) ListAllUnfinishedInferences(ctx context.Context) ([]ListAllUnfinishedInferencesRow, error) {
	rows, err := q.db.Query(ctx, listAllUnfinishedInferences)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAllUnfinishedInferencesRow
	for rows.Next() {
		var i ListAllUnfinishedInferencesRow
		if err := rows.Scan(
			&i.InferenceID,
			&i.Prompt,
			&i.NegPrompt,
			&i.FinishedAt,
			&i.ID,
			&i.Domain,
			&i.Name,
			&i.ParentID,
			&i.Status,
			&i.SampleID,
			&i.TrainID,
			&i.UpdatedAt,
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

const listAllUnfinishedSample = `-- name: ListAllUnfinishedSample :many
SELECT id, model_id, finished_at, created_at FROM samples
WHERE finished_at IS NULL
`

func (q *Queries) ListAllUnfinishedSample(ctx context.Context) ([]Sample, error) {
	rows, err := q.db.Query(ctx, listAllUnfinishedSample)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Sample
	for rows.Next() {
		var i Sample
		if err := rows.Scan(
			&i.ID,
			&i.ModelID,
			&i.FinishedAt,
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

const listAllUnfinishedTrain = `-- name: ListAllUnfinishedTrain :many
SELECT trains.id as train_id, m.id, m.domain, m.name, m.parent_id, m.status, m.sample_id, m.train_id, m.updated_at, m.created_at FROM trains
JOIN models m ON m.id = trains.id
WHERE trains.finished_at IS NULL
`

type ListAllUnfinishedTrainRow struct {
	TrainID   string
	ID        string
	Domain    string
	Name      string
	ParentID  pgtype.Text
	Status    string
	SampleID  pgtype.Text
	TrainID_2 pgtype.Text
	UpdatedAt pgtype.Timestamptz
	CreatedAt pgtype.Timestamptz
}

func (q *Queries) ListAllUnfinishedTrain(ctx context.Context) ([]ListAllUnfinishedTrainRow, error) {
	rows, err := q.db.Query(ctx, listAllUnfinishedTrain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAllUnfinishedTrainRow
	for rows.Next() {
		var i ListAllUnfinishedTrainRow
		if err := rows.Scan(
			&i.TrainID,
			&i.ID,
			&i.Domain,
			&i.Name,
			&i.ParentID,
			&i.Status,
			&i.SampleID,
			&i.TrainID_2,
			&i.UpdatedAt,
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

const listFinishedInferences = `-- name: ListFinishedInferences :many
SELECT i.id, model_id, prompt, neg_prompt, finished_at, m.id, domain, name, parent_id, status, sample_id, train_id, updated_at, created_at FROM inferences i
JOIN models m ON i.model_id = m.id
WHERE finished_at IS NOT NULL
ORDER BY finished_at DESC
LIMIT $1 OFFSET $2
`

type ListFinishedInferencesParams struct {
	Limit  int32
	Offset int32
}

type ListFinishedInferencesRow struct {
	ID         string
	ModelID    string
	Prompt     string
	NegPrompt  string
	FinishedAt pgtype.Timestamptz
	ID_2       string
	Domain     string
	Name       string
	ParentID   pgtype.Text
	Status     string
	SampleID   pgtype.Text
	TrainID    pgtype.Text
	UpdatedAt  pgtype.Timestamptz
	CreatedAt  pgtype.Timestamptz
}

func (q *Queries) ListFinishedInferences(ctx context.Context, arg ListFinishedInferencesParams) ([]ListFinishedInferencesRow, error) {
	rows, err := q.db.Query(ctx, listFinishedInferences, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListFinishedInferencesRow
	for rows.Next() {
		var i ListFinishedInferencesRow
		if err := rows.Scan(
			&i.ID,
			&i.ModelID,
			&i.Prompt,
			&i.NegPrompt,
			&i.FinishedAt,
			&i.ID_2,
			&i.Domain,
			&i.Name,
			&i.ParentID,
			&i.Status,
			&i.SampleID,
			&i.TrainID,
			&i.UpdatedAt,
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

const listModelByDomain = `-- name: ListModelByDomain :many
SELECT id, domain, name, parent_id, status, sample_id, train_id, updated_at, created_at
FROM models
WHERE domain = $1
`

func (q *Queries) ListModelByDomain(ctx context.Context, domain string) ([]Model, error) {
	rows, err := q.db.Query(ctx, listModelByDomain, domain)
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
			&i.ParentID,
			&i.Status,
			&i.SampleID,
			&i.TrainID,
			&i.UpdatedAt,
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

const listModels = `-- name: ListModels :many
SELECT id, domain, name, parent_id, status, sample_id, train_id, updated_at, created_at FROM models
WHERE id = ANY($1::text[])
`

func (q *Queries) ListModels(ctx context.Context, dollar_1 []string) ([]Model, error) {
	rows, err := q.db.Query(ctx, listModels, dollar_1)
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
			&i.ParentID,
			&i.Status,
			&i.SampleID,
			&i.TrainID,
			&i.UpdatedAt,
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

const listTasks = `-- name: ListTasks :many
SELECT task_id, task_type, content, status, estimate, updated_at, created_at FROM tasks
ORDER BY created_at DESC
`

func (q *Queries) ListTasks(ctx context.Context) ([]Task, error) {
	rows, err := q.db.Query(ctx, listTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.TaskID,
			&i.TaskType,
			&i.Content,
			&i.Status,
			&i.Estimate,
			&i.UpdatedAt,
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

const updateInferenceFinished = `-- name: UpdateInferenceFinished :exec
UPDATE inferences
SET finished_at = now()
WHERE id = $1
`

func (q *Queries) UpdateInferenceFinished(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, updateInferenceFinished, id)
	return err
}

const updateModelStatus = `-- name: UpdateModelStatus :exec
UPDATE models
SET status = $2
WHERE id = $1
`

type UpdateModelStatusParams struct {
	ID     string
	Status string
}

func (q *Queries) UpdateModelStatus(ctx context.Context, arg UpdateModelStatusParams) error {
	_, err := q.db.Exec(ctx, updateModelStatus, arg.ID, arg.Status)
	return err
}

const updateSampleFinished = `-- name: UpdateSampleFinished :exec
UPDATE samples
SET finished_at = now()
WHERE id = $1
`

func (q *Queries) UpdateSampleFinished(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, updateSampleFinished, id)
	return err
}

const updateTaskStatus = `-- name: UpdateTaskStatus :exec
UPDATE tasks
SET status = $2,
    estimate = $3,
    updated_at = now()
WHERE task_id = $1
`

type UpdateTaskStatusParams struct {
	TaskID   int32
	Status   string
	Estimate int64
}

func (q *Queries) UpdateTaskStatus(ctx context.Context, arg UpdateTaskStatusParams) error {
	_, err := q.db.Exec(ctx, updateTaskStatus, arg.TaskID, arg.Status, arg.Estimate)
	return err
}

const updateTrainFinished = `-- name: UpdateTrainFinished :exec
UPDATE trains
SET finished_at = now()
WHERE id = $1
`

func (q *Queries) UpdateTrainFinished(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, updateTrainFinished, id)
	return err
}
