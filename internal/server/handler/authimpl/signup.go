package authimpl

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/errors"
	"github.com/thesis-bkn/hfsd/internal/server/transport"
)

// sign up route
type SignupRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
}

type SignupResponse struct {
	Token string `json:"token"`
}

// Signup implements AuthHandler.
func (a *authHandler) Signup(c echo.Context) error {
	var (
		signupRequest SignupRequest
		user          entity.User
	)
	err := c.Bind(&signupRequest)
	if err != nil {
		return errors.ErrBadRequest
	}

	if err := a.validate.Struct(signupRequest); err != nil {
		return errors.ErrBadRequest
	}

	user.ID = uuid.NewString()
	user.Name = signupRequest.Username
	user.Email = signupRequest.Email
	user.Activated = false
	user.Password, err = hashPassword(signupRequest.Password)
	if err != nil {
		return tracerr.Wrap(err)
	}

	if err = a.userRepo.CreateUser(c.Request().Context(), &user); err != nil {
		return errors.ErrBadRequest
	}

	token, err := newToken(&user, a.cfg.Authenticate.JwtSecret)
	if err != nil {
		return tracerr.Wrap(err)
	}

	cookie := new(http.Cookie)
	cookie.Name = transport.ContextKeyCookieToken.String()
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 24 * 14).UTC()
	cookie.HttpOnly = true

	c.SetCookie(cookie)
	c.JSON(http.StatusOK, SignupResponse{Token: token})

	return nil
}
