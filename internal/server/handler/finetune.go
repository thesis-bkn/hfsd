package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/splode/fname"
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/errors"
	"github.com/thesis-bkn/hfsd/internal/utils"
)

type FinetuneModelHandler struct {
	cc       chan<- entity.Task
	validate *validator.Validate
	client   database.Client
	cfg      *config.Config
	nameRng  *fname.Generator
}

func NewFinetuneModelHandler(
	taskqueue chan<- entity.Task,
	validate *validator.Validate,
	client database.Client,
	nameRng *fname.Generator,
	cfg *config.Config,
) *FinetuneModelHandler {
	return &FinetuneModelHandler{
		validate: validate,
		cc:       taskqueue,
		client:   client,
		nameRng:  nameRng,
		cfg:      cfg,
	}
}

func (h *FinetuneModelHandler) SubmitSampleTask(c echo.Context) error {
	var req struct {
		ModelID string `param:"modelID" validate:"required"`
	}

	var res struct {
		Name string `json:"name"`
	}

	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest
	}

	if err := h.validate.Struct(req); err != nil {
		return tracerr.Wrap(err)
	}

	sourceModel, err := h.client.Query().GetModelByID(c.Request().Context(), req.ModelID)
	if err != nil {
		c.Error(errors.ErrNotFound)
		return tracerr.Wrap(err)
	}

	sourceModelAgg := entity.NewModelFromDB(&sourceModel)

	modelName, err := h.nameRng.Generate()
	if err != nil {
		return tracerr.Wrap(err)
	}

	sample, err := entity.NewSample(sourceModelAgg, true)
	if err != nil {
		return tracerr.Wrap(err)
	}

	res.Name = modelName

	sample.Model().Sampling()
	if err := h.client.Query().InsertModel(c.Request().Context(), sample.Model().Insertion()); err != nil {
		return tracerr.Wrap(err)
	}

	if err := h.client.Query().InsertSample(c.Request().Context(), sample.Insertion()); err != nil {
		return tracerr.Wrap(err)
	}

	// worker do sample task
	h.cc <- sample

	c.JSON(http.StatusOK, res)

	return nil
}

type fineTuneTaskRequest struct {
	ModelID string             `json:"modelID" validate:"required"`
	Items   []fineTuneTaskItem `json:"items"   validate:"gte=0"`
}

type fineTuneTaskItem struct {
	Order  int  `json:"order"`
	Option bool `json:"option"`
}

func (h *FinetuneModelHandler) SubmitFinetuneTask(c echo.Context) error {
	var req fineTuneTaskRequest
	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest
	}

	if err := h.validate.Struct(req); err != nil {
		return tracerr.Wrap(err)
	}

	modelDB, err := h.client.Query().GetModelByID(c.Request().Context(), req.ModelID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	model := entity.NewModelFromDB(&modelDB)
	sample, err := entity.NewSample(model, false)
	if err != nil {
		return tracerr.Wrap(err)
	}

	train := entity.NewTrain(sample)

	// -- Write to json file
	sort.Slice(req.Items, func(i, j int) bool {
		return req.Items[i].Order < req.Items[j].Order
	})
	ratings := utils.Map(req.Items, func(e fineTuneTaskItem) int {
		if e.Option {
			return 0
		}
		return -1
	})
	if err := h.writeJsonFile(model, ratings); err != nil {
		return tracerr.Wrap(err)
	}

	h.cc <- train

	if err := h.client.Query().InsertTrain(c.Request().Context(), train.Insertion()); err != nil {
		return tracerr.Wrap(err)
	}

	modelStatus := sample.Model().Training()
	if err := h.client.
		Query().
		UpdateModelStatus(
			c.Request().Context(),
			database.UpdateModelStatusParams{
				ID:     train.Model().ID(),
				Status: modelStatus.String(),
			}); err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (h *FinetuneModelHandler) writeJsonFile(model *entity.Model, ratings []int) error {
	ratingGroups := [][]int{}
	for i := 0; i < len(ratings); i += 7 {
		ratingGroups = append(ratingGroups, []int{})
		for j := 0; j < 7; j++ {
			ratingGroups[i/7] = append(ratingGroups[i/7], ratings[i+j])
		}
	}

	// Open file for writing
	if err := os.Mkdir(fmt.Sprintf("./data/assets/samples/%s/json", model.SampleID()), os.ModePerm); err != nil {
		return tracerr.Wrap(err)
	}

	file, err := os.Create(fmt.Sprintf("./data/assets/samples/%s/json/data.json", model.SampleID()))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return tracerr.Wrap(err)
	}
	defer file.Close()

	// Create a JSON encoder and write the nested array to the file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(ratingGroups); err != nil {
		fmt.Println("Error encoding JSON to file:", err)
		return tracerr.Wrap(err)
	}

	return nil
}
