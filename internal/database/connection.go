package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/ztrue/tracerr"
)

type Client interface {
	Query() *Queries
	Conn() *pgx.Conn
}

type client struct {
	queries *Queries
	conn    *pgx.Conn
}

type QuerierT[T any] interface {
	New(db pgx.Conn) *T
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

	return &client{
		queries: queries,
		conn:    conn,
	}, nil
}

func (c *client) Query() *Queries {
	return c.queries
}

func (c *client) Conn() *pgx.Conn {
	return c.conn
}
