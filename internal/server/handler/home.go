package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/templates"
)

type HomeHandler interface {
	Home(c echo.Context) error
}

type homeHandler struct{}

func NewHomeHandler() HomeHandler {
	return &homeHandler{}
}

// Home implements HomeHandler.
func (*homeHandler) Home(c echo.Context) error {
	return render(c, templates.Home())
}
