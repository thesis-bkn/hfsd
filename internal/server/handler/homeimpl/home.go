package homeimpl

import (
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/server/handler"
	"github.com/thesis-bkn/hfsd/templates"
)

type homeHandler struct{}

func NewHomeHandler() handler.HomeHandler {
	return &homeHandler{}
}

// Home implements HomeHandler.
func (*homeHandler) Home(c echo.Context) error {
	return templates.Home().Render(c.Request().Context(), c.Response().Writer)
}
