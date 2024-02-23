package handler

import "github.com/labstack/echo/v4"

type ModelHandler interface {
	Search(c echo.Context) error
}
