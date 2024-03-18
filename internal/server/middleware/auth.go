package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/server/transport"
)

func Authenticate(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			jwtToken, ok := c.Get(transport.ContextKeyCookieToken.String()).(string)
			if !ok {
				c.Redirect(http.StatusMovedPermanently, "/auth/login")
				return nil
			}
			claims := &entity.ProfileClaim{}

			if _, err := jwt.ParseWithClaims(jwtToken, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(cfg.Authenticate.JwtSecret), nil
			}); err != nil {
				c.Redirect(http.StatusMovedPermanently, "/auth/login")
				return nil
			}

			c.Set(transport.ContextAuthenticatedUID.String(), claims.UserID)

			return next(c)
		}
	}
}
