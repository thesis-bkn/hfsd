package view

import (
	"github.com/labstack/echo/v4"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
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
	modelTrains, err := f.client.
		Query().
		ListAllUnfinishedTrain(c.Request().Context())
	if err != nil {
		return tracerr.Wrap(err)
	}

	modelInfs, err := f.client.
		Query().ListAllUnfinishedInferences(c.Request().Context())
	if err != nil {
		return tracerr.Wrap(err)
	}

	trains, err := utils.MapErr(modelTrains, modelTrainsToTrain)
	infs, err := utils.MapErr(modelInfs, modelInfsToInfs)

	var tasks []templates.Task
	tasks = utils.Map(trains, trainToTask)
	tasks = append(tasks, utils.Map(infs, infToTask)...)

	return templates.
		FactoryView(f.cfg.BucketEpt(), tasks).
		Render(c.Request().Context(), c.Response().Writer)
}

func trainToTask(t *entity.Train) templates.Task {
	return templates.Task{
		ID:       t.ID(),
		Type:     templates.Train,
		ImageUrl: t.GetSample().ViewImage(),
		Status:   0,
	}
}

func infToTask(i *entity.Inference) templates.Task {
	return templates.Task{
		ID:       i.ID(),
		Type:     templates.Inference,
		ImageUrl: i.ViewImage(),
		Status:   0,
	}
}

func modelInfsToInfs(mi database.ListAllUnfinishedInferencesRow) (*entity.Inference, error) {
	var parentID *string
	if mi.ParentID.Valid {
		parentID = &mi.ParentID.String
	}
	m, err := entity.NewModel(
		mi.Domain,
		parentID,
		mi.Status,
		mi.SampleID.String,
		mi.TrainID.String,
		mi.ID,
	)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return entity.NewInferenceFromModel(m,
		database.Inference{
			ID:         mi.ID,
			ModelID:    m.ID(),
			Prompt:     mi.Prompt,
			NegPrompt:  mi.NegPrompt,
			FinishedAt: mi.FinishedAt,
		}), nil
}

func modelTrainsToTrain(mt database.ListAllUnfinishedTrainRow) (*entity.Train, error) {
	var parentID *string
	if mt.ParentID.Valid {
		parentID = &mt.ParentID.String
	}
	m, err := entity.NewModel(
		mt.Domain,
		parentID,
		mt.Status,
		mt.SampleID.String,
		mt.TrainID,
		mt.ID,
	)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return entity.NewTrainFromModel(m)
}
