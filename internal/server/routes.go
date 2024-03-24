package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/s3"
	"github.com/thesis-bkn/hfsd/internal/server/handler"
	appmw "github.com/thesis-bkn/hfsd/internal/server/middleware"
	"github.com/thesis-bkn/hfsd/internal/server/view"
	"github.com/thesis-bkn/hfsd/templates"
)

func registerRoutes(cfg *config.Config, client database.Client) *echo.Echo {
	e := echo.New()
	e.Server.Addr = fmt.Sprintf(":%d", cfg.Port)
	e.Server.IdleTimeout = time.Minute
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 30 * time.Second

	// Helper
	validate := validator.New(validator.WithRequiredStructEnabled())
	s3Client := s3.NewS3Client(cfg)

	// Handler
	inferenceHandler := handler.NewInferenceHandler(cfg, s3Client, client, validate)

	// View
	homeView := view.NewHomeView()
	infView := view.NewInferenceView()
	finetuneView := view.NewFinetuneView()
	factoryView := view.NewFactoryView(client, cfg)

	// Global Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(appmw.PopulateRequestContext())
	e.Use(appmw.PopulateCookieContext())

	// Static file for js scripts (htmx,...)
	fileServer := http.FileServer(http.FS(templates.Files))
	assetFile := http.FileServer(http.FS(templates.Assets))
	e.GET("/static/*", echo.WrapHandler(fileServer))
	e.GET("/asset/*", echo.WrapHandler(assetFile))

	e.GET("/", homeView.Home)
	// APIs
	apiEpts := e.Group("/api")
	apiEpts.POST("/inference", inferenceHandler.SubmitInferenceTask)

	// Views ------
	e.GET("/inference", infView.View)
	e.GET("/factory", factoryView.View)
	e.GET("/finetune", finetuneView.FinetuneView)
	e.GET("/finetune/:clientID", finetuneView.FinetuneModelView)

	return e
}
