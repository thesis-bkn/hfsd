package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/ztrue/tracerr"
)

type Client interface {
	Query() *Queries
	Conn() *pgxpool.Pool
}

type client struct {
	queries *Queries
	conn    *pgxpool.Pool
}

func NewClient(cfg *config.Config) (Client, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	queries := New(pool)

	return &client{
		queries: queries,
		conn:    pool,
	}, nil
}

func (c *client) Query() *Queries {
	return c.queries
}

func (c *client) Conn() *pgxpool.Pool {
	return c.conn
}
