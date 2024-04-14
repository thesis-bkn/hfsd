package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/splode/fname"
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/errors"
)

type FinetuneModelHandler struct {
	validate *validator.Validate
	client   database.Client
	cfg      *config.Config
	nameRng  *fname.Generator
}

func NewFinetuneModelHandler(
	validate *validator.Validate,
	client database.Client,
	nameRng *fname.Generator,
	cfg *config.Config,
) *FinetuneModelHandler {
	return &FinetuneModelHandler{
		validate: validate,
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

	sourceModel, err := h.client.Query().GetModel(c.Request().Context(), req.ModelID)
	if err != nil {
		c.Error(errors.ErrNotFound)
		return tracerr.Wrap(err)
	}

	newModelID, err := gonanoid.Generate(modelIDsAlphabet, 5)
	if err != nil {
		return tracerr.Wrap(err)
	}

	modelName, err := h.nameRng.Generate()
	if err != nil {
		return tracerr.Wrap(err)
	}

	newModel := &database.Model{
		ID:     newModelID,
		Domain: sourceModel.Domain,
		Name:   modelName,
		Parent: req.ModelID,
	}

	if err := h.client.Query().InsertPendingModel(c.Request().Context(), database.InsertPendingModelParams{
		ID:     newModel.ID,
		Domain: newModel.Domain,
		Name:   newModel.Name,
		Parent: newModel.Parent,
	}); err != nil {
		return tracerr.Wrap(err)
	}

	if err := h.client.Query().InsertSampleTask(c.Request().Context(), database.InsertSampleTaskParams{
		SourceModelID: newModel.Parent,
		OutputModelID: pgtype.Text{
			String: newModel.ID,
			Valid:  true,
		},
	}); err != nil {
		return tracerr.Wrap(err)
	}

	res.Name = newModel.Name
	c.JSON(http.StatusOK, res)

	return nil
}

func (h *FinetuneModelHandler) SubmitFinetuneTask(c echo.Context) error {
	fmt.Println("hello world")

	return nil
}
