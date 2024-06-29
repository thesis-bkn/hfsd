package view

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/utils"
	"github.com/thesis-bkn/hfsd/templates"
	"github.com/ztrue/tracerr"
)

type FactoryView struct {
	client database.Client
	cfg    *config.Config
}

func NewFactoryView(
	client database.Client,
	cfg *config.Config,
) *FactoryView {
	return &FactoryView{
		client: client,
		cfg:    cfg,
	}
}

func (f *FactoryView) View(c echo.Context) error {
	tasks, err := f.client.Query().ListTasks(c.Request().Context())
	if err != nil {
		return tracerr.Wrap(err)
	}

	viewTasks := utils.Map(tasks, func(t database.Task) templates.Task {
		var taskType templates.TaskType
		var status templates.TaskStatus
		var maxV int64 = 100
		var val int64 = 100

		switch t.TaskType {
		case "inference":
			taskType = templates.Inference
		case "train":
			taskType = templates.Train
		case "sample":
			taskType = templates.Sample
		}

		switch t.Status {
		case "pending":
			status = templates.Pending
		case "running":
			status = templates.Processing
		case "finished":
			status = templates.Finished
		}

		if t.Estimate != -1 {
			val = int64(time.Since(t.UpdatedAt.Time).Seconds())
			maxV = t.Estimate
		}

		if t.Estimate == -1 && status == templates.Pending {
			val = 0
		}

		return templates.Task{
			ID:      t.TaskID,
			Type:    taskType,
			Status:  status,
			Content: t.Content,
			Max:     maxV,
			Value:   val,
		}
	})

	return templates.
		FactoryView(f.cfg.BucketEpt(), viewTasks).
		Render(c.Request().Context(), c.Response().Writer)
}
