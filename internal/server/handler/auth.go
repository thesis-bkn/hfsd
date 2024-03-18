package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/errors"
	"github.com/thesis-bkn/hfsd/internal/server/transport"
	"github.com/thesis-bkn/hfsd/templates"
	"github.com/ztrue/tracerr"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	validate *validator.Validate

	cfg *config.Config
}

func NewAuthHandler(
	validate *validator.Validate,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		validate: validate,
		cfg:      cfg,
	}
}

// LoginView implements AuthHandler.
func (a *AuthHandler) LoginView(c echo.Context) error {
	return render(c, templates.Login())
}

type LoginSubmitRequest struct {
	Email    string `validate:"required,email" json:"email"    form:"email"`
	Password string `validate:"required"       json:"password" form:"password"`
}

type LoginSubmitResponse struct {
	Username string `json:"username"`
	UserID   string `json:"id"`
	Token    string `json:"token"`
}

// LoginSubmit implements AuthHandler.
func (a *AuthHandler) LoginSubmit(c echo.Context) error {
	var loginRequest LoginSubmitRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		return errors.ErrBadRequest
	}

	if err := a.validate.Struct(loginRequest); err != nil {
		fmt.Println("err", err)
		return errors.ErrBadRequest
	}

	cookie := new(http.Cookie)
	cookie.Name = transport.ContextKeyCookieToken.String()
	// cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 24 * 14).UTC()
	cookie.HttpOnly = true
	cookie.Path = "/"

	c.SetCookie(cookie)

	c.Response().Status = http.StatusOK

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

// sign up route
type SignupRequest struct {
	Password string `json:"password" validate:"required"       form:"password"`
	Email    string `json:"email"    validate:"required,email" form:"email"`
}

type SignupResponse struct {
	Token string `json:"token"`
}

func (a *AuthHandler) SignupView(c echo.Context) error {
	return render(c, templates.Signup())
}

// Signup implements AuthHandler.
func (a *AuthHandler) SignupSubmit(c echo.Context) error {
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
	user.Email = signupRequest.Email
	user.Activated = false
	user.Password, err = hashPassword(signupRequest.Password)
	if err != nil {
		return tracerr.Wrap(err)
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
	cookie.Path = "/"

	c.SetCookie(cookie)

	c.Response().Status = http.StatusOK

	return nil
}
