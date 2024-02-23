package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/ztrue/tracerr"

	_ "github.com/mattn/go-sqlite3"
)

type Client interface {
	DB() *sqlx.DB
}

type client struct {
	db *sqlx.DB
}

const SCHEMA = `
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    activated INTEGER DEFAULT 0,
    created_at INTEGER DEFAULT (strftime('%s', 'now')) NOT NULL
);
`

func NewClient() Client {
	db, err := sqlx.Connect("sqlite3", "./data/sqlite.db")
	if err != nil {
		tracerr.PrintSourceColor(err)
		log.Fatal(err.Error())
		return nil
	}

	db.MustExec(SCHEMA)

	return &client{db}
}

func (c *client) DB() *sqlx.DB {
	return c.db
}
