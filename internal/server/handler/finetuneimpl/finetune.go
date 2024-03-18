package finetuneimpl

import (
	echo "github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/server/handler"
	"github.com/thesis-bkn/hfsd/templates"
)

type FinetuneHandler struct{}

func NewFineTuneHandler() handler.FineTuneHandler {
	return &FinetuneHandler{}
}

// FinetuneView implements handler.FineTuneHandler.
func (*FinetuneHandler) FinetuneView(c echo.Context) error {
	var ret []*templates.ModelInfo
	for i := 0; i < 100; i++ {
		ret = append(ret, &templates.ModelInfo{
			Name: "hello world",
			ID:   "hello-world",
		})
	}
	return templates.FinetuneView(ret).Render(c.Request().Context(), c.Response().Writer)
}

// FinetuneModelView implements handler.FineTuneHandler.
func (*FinetuneHandler) FinetuneModelView(c echo.Context) error {
	panic("unimplemented")
}
