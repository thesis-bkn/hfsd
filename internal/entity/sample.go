package entity

import (
	"fmt"
	"path"

	"github.com/google/uuid"
	"github.com/thesis-bkn/hfsd/data"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/ztrue/tracerr"
)

type Sample struct {
	model *Model

	id string
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
		ID:      s.id,
		ModelID: s.model.id,
	}
}

func (s *Sample) Model() *Model {
	return s.model
}

func (s *Sample) SaveDir() string {
	return path.Join("./data", data.SamplePath, s.id)
}

func (s *Sample) ViewImage() string {
	files, err := data.Samples.ReadDir(s.id)
	if err != nil {
		fmt.Println("should not err here")
	}

	return path.Join("/data", data.SamplePath, s.id, "images", files[0].Name())
}
