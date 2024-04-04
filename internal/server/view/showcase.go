package view

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/templates"
	"github.com/ztrue/tracerr"
)

type ShowcaseView struct {
	client database.Client
	cfg    *config.Config
}

func NewShowcaseView(
	client database.Client,
	cfg *config.Config,
) *ShowcaseView {
	return &ShowcaseView{
		client: client, cfg: cfg,
	}
}

const LIMIT_PER_PAGE = 50

func (v *ShowcaseView) View(c echo.Context) error {
	query := c.Request().URL.Query()
	pageRaw := query.Get("page")

	var (
		page int64 = 0
		err  error
	)
	if pageRaw != "" {
		page, err = strconv.ParseInt(pageRaw, 0, 64)
		if err != nil {
			return tracerr.Wrap(err)
		}
	}
	infs, err := v.client.
		Query().
		ListInferences(
			c.Request().Context(),
			database.ListInferencesParams{
				Limit:  int32(LIMIT_PER_PAGE),
				Offset: int32(page * LIMIT_PER_PAGE),
			},
		)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return templates.
		ShowcaseView(v.cfg.BucketEpt(), infs, page).
		Render(
			c.Request().Context(),
			c.Response().Writer,
		)
}
