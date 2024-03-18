package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/thesis-bkn/hfsd/internal/client"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/server/handler/finetuneimpl"
	"github.com/thesis-bkn/hfsd/internal/server/handler/homeimpl"
	"github.com/thesis-bkn/hfsd/internal/server/handler/inferenceimpl"
	appmw "github.com/thesis-bkn/hfsd/internal/server/middleware"
	"github.com/thesis-bkn/hfsd/templates"
)

func registerRoutes(cfg *config.Config, client client.Client) *echo.Echo {
	e := echo.New()
	e.Server.Addr = fmt.Sprintf(":%d", cfg.Port)
	e.Server.IdleTimeout = time.Minute
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 30 * time.Second

	// Helper
	// validate := validator.New(validator.WithRequiredStructEnabled())

	// Repo

	// Middleware
	mwAuthenticate := appmw.Authenticate(cfg)

	// Handler
	homeHandler := homeimpl.NewHomeHandler()
	infHandler := inferenceimpl.NewInferenceHandler()
	finetuneHandler := finetuneimpl.NewFineTuneHandler()

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

	e.GET("/", homeHandler.Home)

	// Auth ----------
	// ---------------

	// Inference ------
	e.GET("/inference", infHandler.InferenceView)
	// Finetune  -------
	e.GET("/finetune", mwAuthenticate(finetuneHandler.FinetuneView))
	e.GET("/finetune/:clientID", mwAuthenticate(finetuneHandler.FinetuneModelView))

	return e
}
