package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDb(ctx context.Context, dbDsn string) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, dbDsn)
	if err != nil {
		return nil, err
	}
	return newDatabase(pool), nil
}
