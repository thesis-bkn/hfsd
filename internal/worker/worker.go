package worker

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
)

type Worker struct{}

func NewWorker() *Worker {
	return &Worker{}
}

type taskWrapper struct {
	taskID int32
	data   entity.Task
}

func (w *Worker) Run(
	c <-chan entity.Task,
	client database.Client,
) {
	runningTask := false
	pending := []*taskWrapper{}

	execute := func(x []string) {
		runningTask = true
		defer func() {
			runningTask = false
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

	handler := func(wrapper *taskWrapper) {
		switch task := wrapper.data.(type) {
		case *entity.Sample:
			fmt.Println("receive sampling task")
			execute(fmtSample(task))

		case *entity.Train:
			execute(fmtTrain(task))

		case *entity.Inference:
			fmt.Println("receive inference task")
			execute(fmtInf(task))
		}

		if err := client.Query().
			UpdateTaskStatus(context.Background(), database.UpdateTaskStatusParams{
				TaskID:   wrapper.taskID,
				Status:   "finished",
				Estimate: -1,
			}); err != nil {
			fmt.Println("error update task failed: ", err.Error())
			return
		}
	}

	for {
		select {
		case task := <-c:
			taskID, err := client.Query().
				InsertTask(context.Background(), database.InsertTaskParams{
					TaskType: task.TaskType(),
					Content:  task.TaskContent(),
					Status:   "pending",
					Estimate: -1,
				})
			if err != nil {
				fmt.Println("error update task sampled: ", err.Error())
				return
			}

			pending = append(pending, &taskWrapper{
				taskID: taskID,
				data:   task,
			})

		default:
			if !runningTask && len(pending) != 0 {
				task := pending[0]
				pending = pending[1:]

				if err := client.Query().
					UpdateTaskStatus(context.Background(), database.UpdateTaskStatusParams{
						TaskID:   task.taskID,
						Status:   "running",
						Estimate: int64(task.data.Estimate().Seconds()),
					}); err != nil {
					continue
				}

				go handler(task)
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func fmtSample(task *entity.Sample) []string {
    sampleScript := "./d3po/scripts/sample_inpaint.py"
    if task.Model().Domain() == entity.DomainHuman || task.Model().Domain() == entity.DomainLandscape {
        sampleScript = "./d3po/scripts/sample_outpaint.py"
    }
	res := []string{
		"run",
		"accelerate",
		"launch",
        sampleScript,
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
