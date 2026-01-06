package pkg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

func WithTx(ctx context.Context, pool *pgxpool.Pool, fn func(context.Context) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	err = fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return err // Return original error, but log rollback error if needed
		}
		return err
	}

	return tx.Commit(ctx)
}

func GetTx(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if !ok {
		return nil
	}
	return tx
}

func HasTx(ctx context.Context) bool {
	return GetTx(ctx) != nil
}
