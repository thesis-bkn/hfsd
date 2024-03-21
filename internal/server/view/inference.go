package view

import (
	echo "github.com/labstack/echo/v4"

	"github.com/thesis-bkn/hfsd/templates"
)

type InferenceView struct{}

func NewInferenceView() *InferenceView {
	return &InferenceView{}
}

func (*InferenceView) InferenceView(c echo.Context) error {
	return templates.InferenceView().Render(c.Request().Context(), c.Response().Writer)
}
