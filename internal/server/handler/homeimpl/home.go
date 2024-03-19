package homeimpl

import (
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/templates"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// Home implements HomeHandler.
func (*HomeHandler) Home(c echo.Context) error {
	return templates.Home().Render(c.Request().Context(), c.Response().Writer)
}
