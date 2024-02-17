package handler

import (
	"github.com/labstack/echo/v4"
)

type HomeHandler interface {
	Home(c echo.Context) error
}
