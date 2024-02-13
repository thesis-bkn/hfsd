package repo

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/ztrue/tracerr"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userRepo struct {
	client database.Client
}

func NewUserRepo(client database.Client) UserRepo {
	return &userRepo{client}
}

const (
	CREATE_USER_QUERY = `
    insert into users (id, name, password, activated, email) 
    values (:id, :name, :password, :activated, :email)
    `
)

// CreateUser implements UserRepo.
func (u *userRepo) CreateUser(ctx context.Context, user *entity.User) error {
	if _, err := u.client.DB().NamedExecContext(ctx, CREATE_USER_QUERY, user); err != nil {
		return tracerr.Wrap(err.(*pgconn.PgError))
	}

	return nil
}

const (
	GET_BY_USERNAME_QUERY = `
        select * from users 
        where email = $1
    `
)

// GetByUsername implements UserRepo.
func (u *userRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := u.client.DB().GetContext(ctx, &user, GET_BY_USERNAME_QUERY, email); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &user, nil
}
