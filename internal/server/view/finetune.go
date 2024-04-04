package view

import (
	echo "github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/templates"
)

type FinetuneView struct{}

func NewFinetuneView() *FinetuneView {
	return &FinetuneView{}
}

// View implements handler.FineTuneHandler.
func (*FinetuneView) View(c echo.Context) error {
	var ret []*templates.ModelInfo
	for i := 0; i < 100; i++ {
		ret = append(ret, &templates.ModelInfo{
			Name: "hello world",
			ID:   "hello-world",
		})
	}
	return templates.FinetuneView(ret).Render(c.Request().Context(), c.Response().Writer)
}
