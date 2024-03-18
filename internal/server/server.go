package server

import (
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/client"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/ztrue/tracerr"
)

func NewServer() (*echo.Echo, error) {
	cfg := config.LoadConfig()
	client, err := client.NewClient(cfg)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return registerRoutes(cfg, client), nil
}
