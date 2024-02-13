package middleware

import (
	stderr "errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/thesis-bkn/hfsd/internal/errors"
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
			for _, v := range transport.ContextKeyCookieValues() {
				cookie, err := c.Cookie(v.String())
				if err != nil && stderr.Is(err, http.ErrNoCookie) {
					if stderr.Is(err, http.ErrNoCookie) {
						continue
					}
					c.Error(err)
					return errors.ErrBadRequest
				}
				c.Set(v.String(), cookie)
			}

			return next(c)
		}
	}
}
