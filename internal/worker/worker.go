package worker

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/thesis-bkn/hfsd/internal/entity"
)

type Worker struct{}

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) Run(c <-chan interface{}) {
	runningTask := 0
	pending := []interface{}{}

	execute := func(x []string) {
		runningTask++
		defer func() {
			runningTask = max(0, runningTask - 1)
		}()

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

	handler := func(data interface{}) {
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

	for {
		select {
		case data := <-c:
			pending = append(pending, data)
		default:
			if runningTask == 0 && len(pending) != 0 {
				task := pending[0]
				pending = pending[1:]

                go handler(task)
			}

			time.Sleep(1 * time.Second)
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
		"--llog_dir", task.Model().LogDir(),
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
