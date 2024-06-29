package server

import (
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
)

func NewServer(
	taskqueue chan<- entity.Task,
	cfg *config.Config,
	client database.Client,
) (*echo.Echo, error) {
	return registerRoutes(cfg, client, taskqueue), nil
}
