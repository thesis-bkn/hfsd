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
	"github.com/thesis-bkn/hfsd/internal/repo"
	"github.com/thesis-bkn/hfsd/internal/server/handler/authimpl"
	"github.com/thesis-bkn/hfsd/internal/server/handler/homeimpl"
	"github.com/thesis-bkn/hfsd/internal/server/handler/inferenceimpl"
	appmw "github.com/thesis-bkn/hfsd/internal/server/middleware"
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

	// Repo
	userRepo := repo.NewUserRepo(client)

	// Middleware
	mwAuthenticate := appmw.Authenticate(cfg)

	// Handler
	homeHandler := homeimpl.NewHomeHandler()
	authHandler := authimpl.NewAuthHandler(validate, cfg, userRepo)
	infHandler := inferenceimpl.NewInferenceHandler()

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
	authEndpoint := e.Group("/auth")
	authEndpoint.GET("/login", authHandler.LoginView)
	authEndpoint.POST("/login", authHandler.LoginSubmit)

	authEndpoint.GET("/signup", authHandler.SignupView)
	authEndpoint.POST("/signup", authHandler.SignupSubmit)

	authEndpoint.GET("/verify", mwAuthenticate(authHandler.Validate))
	// ---------------

	// Inference ------
	e.GET("/inference", mwAuthenticate(infHandler.InferenceView))

	return e
}
