package cli

import (
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/database/cache"
	"github.com/Ulqiora/Route256Project/internal/database/postgresql"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
)

type PickPointRepository struct {
	db      postgresql.PGXDatabase
	manager *trmpgx.CtxGetter
	cache   cache.Cache
}

func New(db postgresql.PGXDatabase, cache cache.Cache) *PickPointRepository {
	return &PickPointRepository{
		db:      db,
		manager: trmpgx.DefaultCtxGetter,
		cache:   cache,
	}
}

func hashFunction(id string) string {
	return fmt.Sprintf("pickpoint_%s", id)
}
