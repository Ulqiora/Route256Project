package cli

import (
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/database/cache"
	"github.com/Ulqiora/Route256Project/internal/database/postgresql"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
)

type ClientRepository struct {
	db      postgresql.PGXDatabase
	cache   cache.Cache
	manager *trmpgx.CtxGetter
}

func New(db postgresql.PGXDatabase, cache cache.Cache) *ClientRepository {
	return &ClientRepository{
		db:      db,
		cache:   cache,
		manager: trmpgx.DefaultCtxGetter,
	}
}

func hashFunction(id string) string {
	return fmt.Sprintf("pickpoint_%d", id)
}
