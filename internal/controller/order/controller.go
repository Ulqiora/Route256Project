package order

import (
	"context"

	"github.com/jackc/pgx/v4"
	"homework/internal/model"
	"homework/internal/model/order_changers"
	"homework/internal/repository"
)

type Storage interface {
	Create(ctx context.Context, dto repository.OrderDTO) (uint64, error)

	GetByCustomerID(ctx context.Context, id uint64) ([]repository.OrderDTO, error)
	GetByID(ctx context.Context, id uint64) (repository.OrderDTO, error)

	List(ctx context.Context) ([]repository.OrderDTO, error)
	ListReadyToIssued(ctx context.Context) ([]repository.OrderDTO, error)

	UpdateToReceived(ctx context.Context, orderID uint64) error
	UpdateToReturned(ctx context.Context, orderID uint64) error
	Delete(ctx context.Context, ID uint64) error
}

type TransactionManager interface {
	Run(ctx context.Context, options pgx.TxOptions, f func(ctxTX context.Context) error) error
}

type ControllerOrder struct {
	storage  Storage
	changers map[model.TypePacking]order_changers.ChangerOrder
	tm       TransactionManager
}

func New(repository Storage, changers map[model.TypePacking]order_changers.ChangerOrder, tm TransactionManager) *ControllerOrder {
	return &ControllerOrder{
		storage:  repository,
		changers: changers,
		tm:       tm,
	}
}
