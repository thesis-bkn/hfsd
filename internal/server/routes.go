package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/splode/fname"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/s3"
	"github.com/thesis-bkn/hfsd/internal/server/handler"
	appmw "github.com/thesis-bkn/hfsd/internal/server/middleware"
	"github.com/thesis-bkn/hfsd/internal/server/view"
	"github.com/thesis-bkn/hfsd/templates"
)

func registerRoutes(
	cfg *config.Config,
	client database.Client,
	taskqueue chan<- entity.Task,
) *echo.Echo {
	e := echo.New()
	e.Server.Addr = fmt.Sprintf(":%d", cfg.Port)
	e.Server.IdleTimeout = time.Minute
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 30 * time.Second

	// Helper
	validate := validator.New(validator.WithRequiredStructEnabled())
	s3Client := s3.NewS3Client(cfg)
	rng := fname.NewGenerator()

	// Handler
	inferenceHandler := handler.NewInferenceHandler(taskqueue, cfg, s3Client, client, validate)
	finetuneHandler := handler.NewFinetuneModelHandler(taskqueue, validate, client, rng, cfg)
	// View
	homeView := view.NewHomeView()
	infView := view.NewInferenceView()
	finetuneView := view.NewFinetuneView(cfg, validate, client)
	factoryView := view.NewFactoryView(client, cfg)
	showcaseView := view.NewShowcaseView(client, cfg)

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

	// APIs
	apiEpts := e.Group("/api")
	apiEpts.POST("/inference", inferenceHandler.SubmitInferenceTask)
	apiEpts.POST("/finetune/:modelID", finetuneHandler.SubmitSampleTask)
	apiEpts.POST("/feedback/:modelID", finetuneHandler.SubmitFinetuneTask)

	// Views ------
	e.GET("/", homeView.View)
	e.GET("/inference", infView.View)
	e.GET("/factory", factoryView.View)
	e.GET("/finetune/:domain", finetuneView.View)
	e.GET("/feedback/:modelID", finetuneView.FeedBackView)
	e.GET("/showcase", showcaseView.View)

	return e
}
