package cli

import (
	"fmt"

	"github.com/Ulqiora/Route256Project/internal/database/cache"
	"github.com/Ulqiora/Route256Project/internal/database/postgresql"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
)

type Repository struct {
	db      postgresql.PGXDatabase
	manager *trmpgx.CtxGetter
	cache   cache.Cache
}

func New(cache cache.Cache, db postgresql.PGXDatabase) *Repository {
	return &Repository{
		manager: trmpgx.DefaultCtxGetter,
		cache:   cache,
		db:      db,
	}
}

func hashOrder(id string) string {
	return fmt.Sprintf("order_%s", id)
}

func hashStateOrder(state string) string {
	return fmt.Sprintf("state_order_%s", state)
}
