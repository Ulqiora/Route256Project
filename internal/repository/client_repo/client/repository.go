package cli

import (
	"context"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework/internal/database/cache"
	"homework/internal/repository"
)

type ClientRepository struct {
	db     *pgxpool.Pool
	cache  cache.Cache
	getter *trmpgx.CtxGetter
}

func New(db *pgxpool.Pool, cache cache.Cache) *ClientRepository {
	return &ClientRepository{
		db:    db,
		cache: cache,
	}
}
func (repo *ClientRepository) Create(ctx context.Context, client repository.Client) pgtype.UUID {
	repo.getter.DefaultTrOrDB(context.Background(), repo.db)
	return nil
}

func hashFunction(id int) string {
	return fmt.Sprintf("pickpoint_%d", id)
}
