package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thesis-bkn/hfsd/internal/database"
)

type Inference struct {
	model      *Model
	finishedAt *time.Time

	id     string
	prompt string
	neg    string
}

func NewInference(
	m *Model,
	prompt string,
	negPrompt string,
) *Inference {
	return &Inference{
		id:     uuid.NewString(),
		prompt: prompt,
		neg:    negPrompt,
		model:  m,
	}
}

func NewInferenceFromModel(
	m *Model,
	i database.Inference,
) *Inference {
	var f *time.Time
	if i.FinishedAt.Valid {
		f = &i.FinishedAt.Time
	}
	return &Inference{
		id:         i.ID,
		prompt:     i.Prompt,
		neg:        i.NegPrompt,
		model:      m,
		finishedAt: f,
	}
}

func NewInferenceFromJoinedModel(
    r database.ListInferencesRow,
) *Inference {
    model := NewModelFromDB(&database.Model{
    	ID:        r.ModelID,
    	Domain:    r.Domain,
    	Name:      r.Name,
    	ParentID:  r.ParentID,
    	Status:    r.Status,
    	SampleID:  r.SampleID,
    	TrainID:   r.TrainID,
    	UpdatedAt: r.UpdatedAt,
    	CreatedAt: r.CreatedAt,
    })

    return NewInferenceFromModel(model, database.Inference{
    	ID:         r.ID,
    	ModelID:    r.ModelID,
    	Prompt:     r.Prompt,
    	NegPrompt:  r.NegPrompt,
    	FinishedAt: r.FinishedAt,
    })
}

func (i *Inference) ID() string {
	return i.id
}

func (i *Inference) ImagePath() string {
	return fmt.Sprintf("./data/inferences/in_%s.png", i.id)
}

func (i *Inference) ViewImage() string {
	return fmt.Sprintf("/data/inferences/in_%s.png", i.id)
}

func (i *Inference) MaskPath() string {
	return fmt.Sprintf("./data/inferences/ms_%s.png", i.id)
}

func (i *Inference) OutputPath() string {
	return fmt.Sprintf("./data/inferences/ou_%s.png", i.id)
}

func (i *Inference) Prompt() string {
	return i.prompt
}

func (i *Inference) NegPrompt() string {
	return i.neg
}

func (i *Inference) Model() *Model {
	return i.model
}

func (i *Inference) Insertion() database.InsertInferenceParams {
	return database.InsertInferenceParams{
		ID:        i.id,
		ModelID:   i.model.id,
		Prompt:    i.prompt,
		NegPrompt: i.neg,
	}
}
