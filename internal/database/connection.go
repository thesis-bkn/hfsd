package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/ztrue/tracerr"
)

type Client interface {
	Query() *Queries
}

type client struct {
	queries *Queries
}

func NewClient(cfg *config.Config) (Client, error) {
	conn, err := pgx.Connect(
		context.Background(),
		cfg.DatabaseURL,
	)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	queries := New(conn)

	return &client{queries}, nil
}

func (c *client) Query() *Queries {
	return c.queries
}
