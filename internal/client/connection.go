package client

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/thesis-bkn/hfsd/database"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/ztrue/tracerr"
)

type Client interface {
	Query() *database.Queries
}

type client struct {
	queries *database.Queries
}

func NewClient(cfg *config.Config) (Client, error) {
	conn, err := pgx.Connect(
		context.Background(),
		cfg.DatabaseURL,
	)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	queries := database.New(conn)

	return &client{queries}, nil
}

func (c *client) Query() *database.Queries {
	return c.queries
}
