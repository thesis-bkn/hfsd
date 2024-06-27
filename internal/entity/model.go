package entity

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/splode/fname"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/ztrue/tracerr"
)

type Model struct {
	parentID *string
	id       string
	sampleID string
	trainID  string
	name     string
	domain   Domain
	status   ModelStatus
}

// ENUM(init, sampling, sampled, rating, training, trained)
//
//go:generate go-enum --marshal --values
type ModelStatus int

func NewModelFromDB(m *database.Model) *Model {
	var parentID *string
	if m.ParentID.Valid {
		parentID = &m.ParentID.String
	}

	mm, _ := NewModel(
		m.Domain,
		parentID,
		m.Status,
		m.SampleID.String,
		m.TrainID.String,
		m.ID, m.Name,
	)

	return mm
}

func NewModel(
	domain string,
	parentID *string,
	status string,
	sampleID string,
	trainID string,
	opts ...string,
) (*Model, error) {
	modelStatus, err := ParseModelStatus(status)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	modelDomain, err := ParseDomain(domain)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	modelID := uuid.NewString()
	if len(opts) != 0 {
		modelID = opts[0]
	}

	modelName, err := fname.NewGenerator().Generate()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	if len(opts) == 2 {
		modelName = opts[1]
	}

	return &Model{
		id:       modelID,
		parentID: parentID,
		name:     modelName,
		domain:   modelDomain,
		status:   modelStatus,
		sampleID: sampleID,
		trainID:  trainID,
	}, nil
}

var errWrongModelStatus = errors.New("wrong model status to do this action")

func (m *Model) NewChild(sampleID string) (*Model, error) {
	if m.status != ModelStatusTrained {
		return nil, errWrongModelStatus
	}

	model, err := NewModel(
		m.domain.String(),
		&m.id,
		ModelStatusInit.String(),
		sampleID,
		uuid.NewString(),
	)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return model, nil
}

func (m *Model) Sampling() {
	m.status = ModelStatusSampling
}

func (m *Model) Training() ModelStatus {
	m.status = ModelStatusTraining
	return m.status
}

func (m *Model) Status() ModelStatus {
	return m.status
}

func (m *Model) Insertion() database.InsertModelParams {
	return database.InsertModelParams{
		ID:     m.id,
		Domain: m.domain.String(),
		ParentID: pgtype.Text{
			String: *m.parentID,
			Valid:  true,
		},
		Name:   m.name,
		Status: m.status.String(),
		SampleID: pgtype.Text{
			String: m.sampleID,
			Valid:  true,
		},
		TrainID: pgtype.Text{
			String: m.trainID,
			Valid:  true,
		},
	}
}

func (m *Model) ID() string {
	return m.id
}

func (m *Model) SampleID() string {
	return m.sampleID
}

func (m *Model) Domain() Domain {
	return m.domain
}

func (m *Model) IsBase() bool {
	return m.parentID == nil || strings.Contains(*m.parentID, "base")
}

func (m *Model) ResumeFrom() string {
	if m.parentID == nil {
		panic("can not resume for base model")
	}

	return fmt.Sprintf("./data/assets/logs/%s", *m.parentID)
}

func (m *Model) LogDir() string {
	if m.parentID == nil {
		panic("can not resume for base model")
	}

	return fmt.Sprintf("./data/assets/logs/%s", m.id)
}

func (m *Model) JsonPath() string {
	return fmt.Sprintf("./data/assets/samples/%s/json", m.sampleID)
}

func (m *Model) SamplePath() string {
	return fmt.Sprintf("./data/assets/samples/%s", m.sampleID)
}
