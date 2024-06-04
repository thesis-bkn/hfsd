package worker

import (
	"fmt"
	"os/exec"

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
	for data := range w.C {
		switch task := data.(type) {
		case *entity.TaskSample:
			cmd := exec.Command(
				"poetry",
				fmtSample(task)...,
			)
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			// Print the output
			fmt.Println(string(stdout))

		case *entity.TaskTrain:
			cmd := exec.Command(
				"poetry",
				fmtTrain(task)...,
			)
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			// Print the output
			fmt.Println(string(stdout))
		}
	}
}

func fmtSample(task *entity.TaskSample) []string {
	res := []string{
		"run",
		"accelerate",
		"launch",
		"./d3po/scripts/sample_inpaint.py",
		"--image_fn", task.ImageFn,
		"--prompt_fn_inpaint", task.PromptFnInpaint,
	}

	if task.ResumeFrom != nil {
		res = append(res, "--resume_from", *task.ResumeFrom)
	}

	return res
}

func fmtTrain(task *entity.TaskTrain) []string {
	res := []string{
		"run",
		"accelerate",
		"launch",
		"./d3po/scripts/train_inpaint.py",
		"--train_json_path", task.TrainJsonPath,
		"--train_sample_path", task.TrainSamplePath,
	}

	if task.ResumeFrom != nil {
		res = append(res, "--resume_from", *task.ResumeFrom)
	}

	return res
}
