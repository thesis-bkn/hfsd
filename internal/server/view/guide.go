package view

import (
	echo "github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/templates"
)

type UserGuideView struct{}

func NewUserGuideView() *UserGuideView {
	return &UserGuideView{}
}

// View implements handler.FineTuneHandler.
func (v *UserGuideView) View(c echo.Context) error {
	return templates.UserGuideView().Render(c.Request().Context(), c.Response().Writer)
}
