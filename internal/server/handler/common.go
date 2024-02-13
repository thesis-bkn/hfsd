package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, render templ.Component) error {
	return render.Render(c.Request().Context(), c.Response().Writer)
}
