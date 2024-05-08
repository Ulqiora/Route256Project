package cli

import (
	"fmt"

	"homework/internal/database/cache"
	"homework/internal/database/transaction_manager"
)

type PickPointRepository struct {
	db    *transaction_manager.TransactionManager
	cache cache.Cache
}

func New(db *transaction_manager.TransactionManager, cache cache.Cache) *PickPointRepository {
	return &PickPointRepository{
		db:    db,
		cache: cache,
	}
}

func hashFunction(id int) string {
	return fmt.Sprintf("pickpoint_%d", id)
}
