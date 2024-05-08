package grpc_service

import (
	"context"

	"homework/internal/database/cache/in_memory_cache"
	"homework/internal/database/transaction_manager"
	OrderRepository "homework/internal/repository/client_repo/order"
	PickpointRepository "homework/internal/repository/client_repo/pickpoint"
)

type Repositories struct {
	Order     *OrderRepository.Repository
	Pickpoint *PickpointRepository.PickPointRepository
}

func ConfigureRepositories(ctx context.Context, txManager *transaction_manager.TransactionManager) (Repositories, error) {
	inMemoryCache := in_memory_cache.New()
	inMemoryCache.Run()
	return Repositories{
		Order:     OrderRepository.New(txManager, inMemoryCache),
		Pickpoint: PickpointRepository.New(txManager, inMemoryCache),
	}, nil
}
