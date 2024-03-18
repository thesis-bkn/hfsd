package handler

import "github.com/labstack/echo/v4"

type FineTuneHandler interface {
	FinetuneView(c echo.Context) error
	FinetuneModelView(c echo.Context) error
}
