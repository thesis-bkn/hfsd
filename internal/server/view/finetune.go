package view

import (
	echo "github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/templates"
)

type FinetuneView struct{}

func NewFinetuneView() *FinetuneView {
	return &FinetuneView{}
}

// View implements handler.FineTuneHandler.
func (*FinetuneView) View(c echo.Context) error {
	models := []database.Model{
		{
			ID:   "base",
			Name: "base",
		},
		{
			ID:     "baobui",
			Name:   "baobuidepchai",
			Parent: "base",
		},
		{
			ID:     "baobui2",
			Name:   "son of baobui",
			Parent: "baobui",
		},
		{
			ID:     "baobui3",
			Name:   "second son of baobui",
			Parent: "baobui",
		},
	}
	return templates.
		FinetuneView(models).
		Render(
			c.Request().Context(),
			c.Response().Writer,
		)
}
