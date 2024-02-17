package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/thesis-bkn/hfsd/internal/server/transport"
)

func PopulateRequestContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, v := range transport.ContextKeyRequestValues() {
				c.Set(v.String(), c.Request().Header.Get(v.String()))
			}

			return next(c)
		}
	}
}

func PopulateCookieContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, cookie := range c.Cookies() {
				c.Set(cookie.Name, cookie.Value)
			}

			return next(c)
		}
	}
}
