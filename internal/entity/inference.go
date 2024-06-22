package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thesis-bkn/hfsd/internal/database"
)

type Inference struct {
	id         string
	prompt     string
	neg        string
	finishedAt *time.Time

	model *Model
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
