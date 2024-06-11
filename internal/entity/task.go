package entity

import (
	"fmt"
	"path"
	"time"
)

const (
	data_dir = "./data"

	inf_dir         = "inferences"
	inf_mask_path   = "image/mask.png"
	inf_image_path  = "images/input.png"
	inf_output_path = "images/output.png"
)

type TaskSample struct {
	SaveDir    time.Time
	ResumeFrom *time.Time
	ImageFn    string
	PromptFn   string
}

type TaskTrain struct {
	LogDir          time.Time
	ResumeFrom      *time.Time
	TrainSamplePath string
	TrainJsonPath   string
}

type TaskInference struct {
	Output     time.Time
	ResumeFrom *time.Time
	Prompt     string
	NegPrompt  string
}

func (t *TaskInference) OutputPath() string {
	return path.Join(
		data_dir,
		inf_dir,
		toUnixS(t.Output),
		inf_output_path,
	)
}

func (t *TaskInference) ImagePath() string {
	return path.Join(
		data_dir,
		inf_dir,
		toUnixS(t.Output),
		inf_image_path,
	)
}

func (t *TaskInference) MaskPath() string {
	return path.Join(
		data_dir,
		inf_dir,
		toUnixS(t.Output),
		inf_mask_path,
	)
}

func toUnixS(t time.Time) string {
	return fmt.Sprintf("%d", t.Unix())
}
