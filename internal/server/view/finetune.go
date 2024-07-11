package view

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	echo "github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/errors"
	"github.com/thesis-bkn/hfsd/internal/utils"
	"github.com/thesis-bkn/hfsd/templates"
)

type FinetuneView struct {
	cfg      *config.Config
	validate *validator.Validate
	client   database.Client
}

func NewFinetuneView(
	cfg *config.Config,
	validate *validator.Validate,
	client database.Client,
) *FinetuneView {
	return &FinetuneView{
		cfg:      cfg,
		validate: validate,
		client:   client,
	}
}

// View implements handler.FineTuneHandler.
func (v *FinetuneView) View(c echo.Context) error {
	var req struct {
		Domain entity.Domain `param:"domain"`
	}

	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest
	}

	if err := v.validate.Struct(req); err != nil {
		return tracerr.Wrap(err)
	}

	models, err := v.client.Query().ListModelByDomain(c.Request().Context(), req.Domain.String())
	if err != nil {
		return tracerr.Wrap(err)
	}

	viewModels := utils.Map(
		models,
		func(m database.Model) templates.ModelNode {
			node := templates.ModelNode{
				ID:   m.ID,
				Name: m.Name,
			}

			if m.Status == "trained" {
				node.Status = templates.Finetuned
			}

			if m.Status == "training" {
				node.Status = templates.Training
			}

			if m.Status == "sampling" {
				node.Status = templates.Sampling
			}

			if m.Status == "rating" {
				node.Status = templates.Rating
			}

			if m.ParentID.Valid {
				node.Parent = &m.ParentID.String
			}

			return node
		})

	return templates.
		FinetuneView(viewModels, req.Domain).
		Render(
			c.Request().Context(),
			c.Response().Writer,
		)
}

type FeedbackEntry struct {
	Name string
}

func (e *FeedbackEntry) Group() int {
	return e.Index() / 7
}

func (e *FeedbackEntry) Order() int {
	return e.Index()
}

func (e *FeedbackEntry) Index() int {
	// Extract the file extension
	ext := filepath.Ext(e.Name)

	// Remove the extension from the filename
	name := strings.TrimSuffix(e.Name, ext)

	// Convert the name to an integer
	num, err := strconv.Atoi(name)
	if err != nil {
		panic(fmt.Sprint("Error converting filename to integer:", err))
	}

	return num
}

func (v *FinetuneView) FeedBackView(c echo.Context) error {
	var req struct {
		ModelID string `param:"modelID"`
	}

	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest
	}

	if err := v.validate.Struct(req); err != nil {
		return tracerr.Wrap(err)
	}

	dbModel, err := v.client.Query().GetModelByID(c.Request().Context(), req.ModelID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	model := entity.NewModelFromDB(&dbModel)
	sample, err := entity.NewSample(model, false)
	if err != nil {
		return tracerr.Wrap(err)
	}

	assets := []templates.FeedbackAsset{}
	files := sample.SampleImages()

	for _, fileName := range files {
		entry := FeedbackEntry{
			Name: fileName,
		}

		assets = append(assets, templates.FeedbackAsset{
			ImageUrl: strings.Join([]string{v.cfg.EndpointUrl, v.cfg.Bucket, sample.ID(), fileName}, "/"),
			Group:    entry.Group(),
			Order:    entry.Order(),
		})
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].Order < assets[j].Order
	})

	return templates.
		FeedBackView(req.ModelID, assets).
		Render(
			c.Request().Context(),
			c.Response().Writer,
		)
}
