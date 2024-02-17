package inferenceimpl

import (
	echo "github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/server/handler"
	"github.com/thesis-bkn/hfsd/templates"
)

type inferenceHandler struct{}

func NewInferenceHandler() handler.InferenceHandler {
	return &inferenceHandler{}
}

// InferenceView implements handler.InferenceHandler.
func (*inferenceHandler) InferenceView(c echo.Context) error {
	return templates.InferenceView().Render(c.Request().Context(), c.Response().Writer)
}

// InferenceSearch implements handler.InferenceHandler.
func (*inferenceHandler) InferenceSearch(c echo.Context) error {
	panic("unimplemented")
}
