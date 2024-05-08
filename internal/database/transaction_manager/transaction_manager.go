package transaction_manager

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"homework/internal/database/postgresql"
)

const Key = "transaction"

type TransactionManager struct {
	pool postgresql.PGXDatabase
}

type QueryEngine interface {
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

func New(pool postgresql.PGXDatabase) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

func (t *TransactionManager) Run(ctx context.Context, options pgx.TxOptions, f func(ctxTX context.Context) error) error {
	tx, err := t.pool.GetPool(ctx).BeginTx(ctx, options)
	if err != nil {
		return err
	}

	if err := f(context.WithValue(ctx, Key, tx)); err != nil {
		errRollback := tx.Rollback(ctx)
		return fmt.Errorf("%s:%s", err, errRollback)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s:%s", err, tx.Rollback(ctx))
	}
	return nil
}

func (t *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(Key).(QueryEngine)
	if ok {
		return tx
	}
	return t.pool
}
