package server

import (
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
)

func NewServer() *echo.Echo {
	cfg := config.LoadConfig()
	client := database.NewClient()

	return registerRoutes(cfg, client)
}
