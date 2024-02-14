package authimpl

import (
	"github.com/go-playground/validator/v10"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/repo"
	"github.com/thesis-bkn/hfsd/internal/server/handler"
)

type authHandler struct {
	validate *validator.Validate

	cfg *config.Config

	userRepo repo.UserRepo
}

func NewAuthHandler(
	validate *validator.Validate,
	cfg *config.Config,
	userRepo repo.UserRepo,
) handler.AuthHandler {
	return &authHandler{
		validate: validate,
		userRepo: userRepo,
		cfg:      cfg,
	}
}
