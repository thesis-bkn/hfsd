package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/thesis-bkn/hfsd/internal/config"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Client interface {
	DB() *sqlx.DB
}

type client struct {
	db *sqlx.DB
}

func NewClient(cfg *config.Config) Client {
	db := sqlx.MustConnect("pgx", cfg.URI)

	return &client{db}
}

func (c *client) DB() *sqlx.DB {
	return c.db
}
