package handler

import (
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	LoginView(c echo.Context) error
	LoginSubmit(c echo.Context) error
	SignupView(c echo.Context) error
	SignupSubmit(c echo.Context) error
	Validate(c echo.Context) error
}
