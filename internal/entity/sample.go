package entity

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/ztrue/tracerr"
)

type Sample struct {
	id string

	model *Model
}

func NewSample(
	sourceModel *Model,
	init bool,
) (*Sample, error) {
	if !init {
		return retrieve(sourceModel), nil
	}
	sampleID := uuid.NewString()
	newModel, err := sourceModel.NewChild(sampleID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &Sample{
		id:    sampleID,
		model: newModel,
	}, nil
}

func retrieve(model *Model) *Sample {
	return &Sample{
		id:    model.sampleID,
		model: model,
	}
}

func (s *Sample) Insertion() database.InsertSampleParams {
	return database.InsertSampleParams{
		ID:       s.id,
		ModelID:  s.model.id,
		SaveDir:  fmt.Sprintf("./data/%s", s.id),
		ImageFn:  s.model.domain.ImageFn(),
		PromptFn: s.model.domain.PromptFn(),
	}
}

func (s *Sample) Model() *Model {
	return s.model
}

func (s *Sample) SaveDir() string {
	return fmt.Sprintf("./data/%s", s.id)
}

func (s *Sample) ViewImage() string {
	return fmt.Sprintf("/data/%s/images/001.png", s.id)
}
