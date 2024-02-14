package handler

import (
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	LoginView(c echo.Context) error
	LoginSubmit(c echo.Context) error
	Signup(c echo.Context) error
	Validate(c echo.Context) error
}
