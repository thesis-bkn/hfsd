package view

import (
    "strings"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/utils"
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

	modelInfs, err := v.client.
		Query().
		ListFinishedInferences(
			c.Request().Context(),
			database.ListFinishedInferencesParams{
				Limit:  int32(LIMIT_PER_PAGE),
				Offset: int32(page * LIMIT_PER_PAGE),
			},
		)
	if err != nil {
		return tracerr.Wrap(err)
	}

	infs := utils.Map(modelInfs, entity.NewInferenceFromJoinedModel)
	showcaseItems := utils.Map(infs, func(i *entity.Inference) templates.ShowcaseItem {
		return templates.ShowcaseItem{
			InputImagePath:  strings.Join([]string{v.cfg.EndpointUrl, v.cfg.Bucket, i.ID(), "in.jpg"}, "/"),
			OutputImagePath: strings.Join([]string{v.cfg.EndpointUrl, v.cfg.Bucket, i.ID(), "out.jpg"}, "/"),
			Prompt:          i.Prompt(),
		}
	})

	return templates.
		ShowcaseView(v.cfg.BucketEpt(), showcaseItems, page).
		Render(
			c.Request().Context(),
			c.Response().Writer,
		)
}
