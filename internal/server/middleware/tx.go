package middleware

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/ztrue/tracerr"
)

type sqlxTxKey struct{}

// TxFromContext returns a Tx stored inside a context, or nil if there isn't one.
func TxFromContext(ctx context.Context) *sqlx.Tx {
	tx, _ := ctx.Value(sqlxTxKey{}).(*sqlx.Tx)
	return tx
}

func Transaction(client database.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			tx, err := client.DB().BeginTxx(c.Request().Context(), nil)
			if err != nil {
				tracerr.Wrap(err)
			}

			defer func() {
				if v := recover(); v != nil {
					_ = tx.Rollback()
					panic(v)
				}
			}()

			ctx = context.WithValue(ctx, sqlxTxKey{}, tx)
			c.SetRequest(c.Request().WithContext(ctx))
			err = next(c)

			if err != nil {
				if rerr := tx.Rollback(); rerr != nil {
					err = tracerr.New(fmt.Sprintf("rolling back transaction: %v", rerr))
				}
				return tracerr.Wrap(err)
			}

			if err := tx.Commit(); err != nil {
				if rerr := tx.Rollback(); rerr != nil {
					err = tracerr.New(fmt.Sprintf("rolling back transaction: %v", rerr))
				}
				return tracerr.Wrap(err)
			}
			return nil
		}
	}
}
