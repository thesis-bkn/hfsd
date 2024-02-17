package handler

import "github.com/labstack/echo/v4"

type InferenceHandler interface {
	InferenceView(c echo.Context) error
	InferenceSearch(c echo.Context) error
}
