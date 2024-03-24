package view

import (
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
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
	tasks, err := f.client.Query().ListAllTaskWithAsset(c.Request().Context())
	if err != nil {
		return tracerr.Wrap(err)
	}

	return templates.FactoryView(f.cfg, tasks).
		Render(c.Request().Context(), c.Response().Writer)
}
