package entity

import (
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/ztrue/tracerr"
)

type Train struct {
	id string

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
