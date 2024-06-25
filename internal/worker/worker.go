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
		fmt.Println("execute task: ", cmd)
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
		case *entity.Sample:
			fmt.Println("receive sampling task")
			execute(fmtSample(task))

		case *entity.Train:
			execute(fmtTrain(task))

		case *entity.Inference:
			fmt.Println("receive inference task")
			execute(fmtInf(task))
		}
	}
}

func fmtSample(task *entity.Sample) []string {
	res := []string{
		"run",
		"accelerate",
		"launch",
		"./d3po/scripts/sample_inpaint.py",
		"--model_id", task.Model().ID(),
		"--sample_id", task.ID(),
		"--save_dir", task.SaveDir(),
		"--image_fn", task.Model().Domain().ImageFn(),
		"--prompt_fn", task.Model().Domain().PromptFn(),
	}

	if !task.Model().IsBase() {
		res = append(res, "--resume_from", task.Model().ResumeFrom())
	}

	return res
}

func fmtTrain(task *entity.Train) []string {
	res := []string{
		"run",
		"accelerate",
		"launch",
		"./d3po/scripts/train_inpaint.py",
		"--model_id", task.Model().ID(),
		"--log_dir", task.Model().LogDir(),
		"--train_json_path", task.Model().JsonPath(),
		"--train_sample_path", task.Model().SamplePath(),
	}

	if !task.Model().IsBase() {
		res = append(res, "--resume_from", task.Model().ResumeFrom())
	}

	return res
}

func fmtInf(task *entity.Inference) []string {
	res := []string{
		"run",
		"accelerate",
		"launch",
		"./d3po/scripts/inference_inpaint.py",
		"--inference_id", task.ID(),
		"--image_path", task.ImagePath(),
		"--mask_path", task.MaskPath(),
		"--output_path", task.OutputPath(),
		"--prompt", fmt.Sprintf("'%s'", task.Prompt()),
		"--neg_prompt", fmt.Sprintf("'%s'", task.NegPrompt()),
	}

	if !task.Model().IsBase() {
		res = append(res, "--resume_from", task.Model().ResumeFrom())
	}

	return res
}

func toUnixS(t time.Time) string {
	return fmt.Sprintf("%d", t.Unix())
}
