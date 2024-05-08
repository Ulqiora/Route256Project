package cli

import (
	"fmt"

	"homework/internal/database/cache"
	"homework/internal/database/transaction_manager"
)

type Repository struct {
	manager *transaction_manager.TransactionManager
	cache   cache.Cache
}

func New(manager *transaction_manager.TransactionManager, cache cache.Cache) *Repository {
	return &Repository{
		manager: manager,
		cache:   cache,
	}
}

func hashOrder(id uint64) string {
	return fmt.Sprintf("order_%d", id)
}

func hashStateOrder(state string) string {
	return fmt.Sprintf("state_order_%s", state)
}
