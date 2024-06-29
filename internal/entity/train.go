package entity

import (
	"fmt"
	"time"

	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/ztrue/tracerr"
)

type Train struct {
	id         string
	finishedAt time.Time

	sample *Sample
}

func NewTrain(sample *Sample) *Train {
	return &Train{
		sample: sample,
	}
}

func NewTrainFromModel(m *Model) (*Train, error) {
	sample, err := NewSample(m, false)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &Train{
		id:     m.trainID,
		sample: sample,
	}, nil
}

func (t *Train) ID() string {
	return t.sample.model.trainID
}

func (t *Train) Insertion() database.InsertTrainParams {
	return database.InsertTrainParams{
		ID:       t.ID(),
		ModelID:  t.sample.model.id,
		SampleID: t.sample.id,
	}
}

func (t *Train) Model() *Model {
	return t.sample.model
}

func (t *Train) GetSample() *Sample {
	return t.sample
}

func (i *Train) TaskContent() string {
	return fmt.Sprintf("train-model-%s", i.sample.model.id)
}

func (i *Train) TaskType() string {
	return "train"
}

func (i *Train) Estimate() time.Duration {
	return time.Minute * 20
}
