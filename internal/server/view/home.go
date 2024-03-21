package view

import (
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/templates"
)

type HomeView struct{}

func NewHomeView() *HomeView {
	return &HomeView{}
}

// Home implements HomeHandler.
func (*HomeView) Home(c echo.Context) error {
	return templates.Home().Render(c.Request().Context(), c.Response().Writer)
}
