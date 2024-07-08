package view

import (
	echo "github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/templates"
)

type UserGuideView struct {
	cfg *config.Config
}

func NewUserGuideView(cfg *config.Config) *UserGuideView {
	return &UserGuideView{cfg}
}

// View implements handler.FineTuneHandler.
func (v *UserGuideView) View(c echo.Context) error {
	return templates.UserGuideView(v.cfg.DownloadDataUrl).
		Render(c.Request().Context(), c.Response().Writer)
}
