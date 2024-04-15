package view

import (
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	echo "github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/errors"
	"github.com/thesis-bkn/hfsd/templates"
	"github.com/ztrue/tracerr"
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

	models, err := v.client.Query().ListModelsByDomain(c.Request().Context(), req.Domain.String())
	if err != nil {
		return tracerr.Wrap(err)
	}

	return templates.
		FinetuneView(models, req.Domain).
		Render(
			c.Request().Context(),
			c.Response().Writer,
		)
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

	assets, err := v.client.Query().ListFeedbackAssetByModelID(c.Request().Context(), pgtype.Text{
		String: req.ModelID,
		Valid:  true,
	})
	if err != nil {
		return tracerr.Wrap(err)
	}

	return templates.
		FeedBackView(v.cfg.BucketEpt(), req.ModelID, assets).
		Render(
			c.Request().Context(),
			c.Response().Writer,
		)
}
