package worker

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/thesis-bkn/hfsd/internal/entity"
)

type Worker struct {
	C chan interface{}
}

func NewWorker() *Worker {
	return &Worker{
		C: make(chan interface{}),
	}
}

func (w *Worker) Run() {
	execute := func(x []string) {
		cmd := exec.Command(
			"poetry",
			x...,
		)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Print the output
		fmt.Println(string(stdout))
	}

	for data := range w.C {
		switch task := data.(type) {
		case *entity.TaskSample:
			execute(fmtSample(task))

		case *entity.TaskTrain:
			execute(fmtTrain(task))

		case *entity.TaskInference:
			execute(fmtInf(task))
		}
	}
}

func fmtSample(task *entity.TaskSample) []string {
	res := []string{
		"run",
		"accelerate",
		"launch",
		"./d3po/scripts/sample_inpaint.py",
		"--save_dir", toUnixS(task.SaveDir),
		"--image_fn", task.ImageFn,
		"--prompt_fn", task.PromptFn,
	}

	if task.ResumeFrom != nil {
		res = append(res, "--resume_from", toUnixS(task.SaveDir))
	}

	return res
}

func fmtTrain(task *entity.TaskTrain) []string {
	res := []string{
		"run",
		"accelerate",
		"launch",
		"./d3po/scripts/train_inpaint.py",
		"--log_dir", toUnixS(task.LogDir),
		"--train_json_path", task.TrainJsonPath,
		"--train_sample_path", task.TrainSamplePath,
	}

	if task.ResumeFrom != nil {
		res = append(res, "--resume_from", toUnixS(*task.ResumeFrom))
	}

	return res
}

func fmtInf(task *entity.TaskInference) []string {
	res := []string{
		"run",
		"accelerate",
		"launch",
		"./d3po/scripts/inference_inpaint.py",
		"--image_path", task.ImagePath(),
		"--mask_path", task.MaskPath(),
		"--output_path", task.OutputPath(),
		"--prompt", task.Prompt,
		"--neg_prompt", task.NegPrompt,
	}

	if task.ResumeFrom != nil {
		res = append(res, "--resume_from", toUnixS(*task.ResumeFrom))
	}

	return res
}

func toUnixS(t time.Time) string {
	return fmt.Sprintf("%d", t.Unix())
}
