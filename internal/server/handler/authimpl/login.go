package authimpl

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/errors"
	"github.com/thesis-bkn/hfsd/internal/server/transport"
	"github.com/thesis-bkn/hfsd/templates"
)

// LoginView implements AuthHandler.
func (a *authHandler) LoginView(c echo.Context) error {
	return render(c, templates.Login())
}

type LoginSubmitRequest struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required"       json:"password"`
}

type LoginSubmitResponse struct {
	Username string `json:"username"`
	UserID   string `json:"id"`
	Token    string `json:"token"`
}

// LoginSubmit implements AuthHandler.
func (a *authHandler) LoginSubmit(c echo.Context) error {
	var loginRequest LoginSubmitRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		return errors.ErrBadRequest
	}

	if err := a.validate.Struct(loginRequest); err != nil {
		return errors.ErrBadRequest
	}

	user, err := a.userRepo.GetByEmail(c.Request().Context(), loginRequest.Email)
	if err != nil {
		fmt.Println(err.Error())
		return errors.ErrBadRequest
	}

	token, err := newToken(user, a.cfg.Authenticate.JwtSecret)
	if err != nil {
		fmt.Println(err.Error())
		return errors.ErrInternalError
	}

	cookie := new(http.Cookie)
	cookie.Name = transport.ContextKeyCookieToken.String()
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 24 * 14).UTC()
	cookie.HttpOnly = true

	c.SetCookie(cookie)

	c.JSON(http.StatusOK, LoginSubmitResponse{
		Username: user.Name,
		UserID:   user.ID,
		Token:    token,
	})

	return nil
}

// Validate implements AuthHandler.
func (a *authHandler) Validate(c echo.Context) error {
	c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "ok",
	})

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func newToken(user *entity.User, jwtSecrets string) (string, error) {
	claims := &entity.ProfileClaim{
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "server",
			Subject:   user.ID,
			Audience:  []string{},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24 * 14).UTC()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		},
		UserID:    user.ID,
		Activated: user.Activated,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecrets))
}

func render(c echo.Context, render templ.Component) error {
	return render.Render(c.Request().Context(), c.Response().Writer)
}
