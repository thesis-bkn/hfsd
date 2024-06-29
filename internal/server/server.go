package server

import (
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/ztrue/tracerr"
)

func NewServer(taskqueue chan<- interface{}) (*echo.Echo, error) {
	cfg := config.LoadConfig()
	client, err := database.NewClient(cfg)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return registerRoutes(cfg, client, taskqueue), nil
}
