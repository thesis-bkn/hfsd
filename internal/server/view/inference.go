package view

import (
	echo "github.com/labstack/echo/v4"

	"github.com/thesis-bkn/hfsd/templates"
)

type InferenceView struct{}

func NewInferenceView() *InferenceView {
	return &InferenceView{}
}

func (*InferenceView) View(c echo.Context) error {
	model := c.Request().URL.Query().Get("model")
	var modelName *string
	if model != "" {
		modelName = &model
	}
	return templates.
		InferenceView(modelName).
		Render(
			c.Request().Context(),
			c.Response().Writer,
		)
}
