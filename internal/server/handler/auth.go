package handler

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/errors"
	"github.com/thesis-bkn/hfsd/internal/repo"
	"github.com/thesis-bkn/hfsd/internal/server/transport"
	"github.com/thesis-bkn/hfsd/templates"
)

type AuthHandler interface {
	Login(c echo.Context) error
	Signup(c echo.Context) error
	Validate(c echo.Context) error
}

type authHandler struct {
	validate *validator.Validate

	cfg *config.Config

	userRepo repo.UserRepo
}

func NewAuthHandler(
	validate *validator.Validate,
	cfg *config.Config,
	userRepo repo.UserRepo,
) AuthHandler {
	return &authHandler{
		validate: validate,
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type LoginRequest struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required"       json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
	UserID   string `json:"id"`
	Token    string `json:"token"`
}

// Login implements AuthHandler.
func (a *authHandler) Login(c echo.Context) error {
	return render(c, templates.Login())

	// var loginRequest LoginRequest
	// err := c.Bind(&loginRequest)
	// if err != nil {
	// 	return errors.ErrBadRequest
	// }
	//
	// if err := a.validate.Struct(loginRequest); err != nil {
	// 	return errors.ErrBadRequest
	// }
	//
	// user, err := a.userRepo.GetByEmail(c.Request().Context(), loginRequest.Email)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return errors.ErrBadRequest
	// }
	//
	// token, err := newToken(user, a.cfg.Authenticate.JwtSecret)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return errors.ErrInternalError
	// }
	//
	// cookie := new(http.Cookie)
	// cookie.Name = transport.ContextKeyCookieToken.String()
	// cookie.Value = token
	// cookie.Expires = time.Now().Add(time.Hour * 24 * 14).UTC()
	// cookie.HttpOnly = true
	//
	// c.SetCookie(cookie)
	//
	// c.JSON(http.StatusOK, LoginResponse{
	// 	Username: user.Name,
	// 	UserID:   user.ID,
	// 	Token:    token,
	// })
	//
	// return nil
}

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

// Validate implements AuthHandler.
func (*authHandler) Validate(c echo.Context) error {
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
