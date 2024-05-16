package grpc_service

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/database/cache/in_memory_cache"
	"github.com/Ulqiora/Route256Project/internal/database/postgresql"
	ClientRepository "github.com/Ulqiora/Route256Project/internal/repository/client_repo/client"
	OrderRepository "github.com/Ulqiora/Route256Project/internal/repository/client_repo/order"
	PickpointRepository "github.com/Ulqiora/Route256Project/internal/repository/client_repo/pickpoint"
)

type Repositories struct {
	Order     *OrderRepository.Repository
	Pickpoint *PickpointRepository.PickPointRepository
	Client    *ClientRepository.ClientRepository
}

func ConfigureRepositories(_ context.Context, pool postgresql.PGXDatabase) (Repositories, error) {
	inMemoryCache := in_memory_cache.New()
	inMemoryCache.Run()
	return Repositories{
		Order:     OrderRepository.New(inMemoryCache, pool),
		Pickpoint: PickpointRepository.New(pool, inMemoryCache),
		Client:    ClientRepository.New(pool, inMemoryCache),
	}, nil
}
