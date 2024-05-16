package order

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/model/order_changers"
	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/jackc/pgtype"
)

type Storage interface {
	Create(ctx context.Context, dto repository.OrderDTO) (pgtype.UUID, error)

	GetByID(ctx context.Context, id pgtype.UUID) (repository.OrderDTO, error)
	GetByCustomerID(ctx context.Context, id pgtype.UUID) ([]repository.OrderDTO, error)

	List(ctx context.Context) ([]repository.OrderDTO, error)
	ListReadyToIssued(ctx context.Context) ([]repository.OrderDTO, error)

	UpdateToReceived(ctx context.Context, orderID pgtype.UUID) error
	UpdateToReturned(ctx context.Context, orderID pgtype.UUID) error

	Delete(ctx context.Context, orderID pgtype.UUID) error
}

type Manager interface {
	// Do processes a transaction inside a closure.
	Do(context.Context, func(ctx context.Context) error) error
	// DoWithSettings processes a transaction inside a closure with custom trm.Settings.
	DoWithSettings(context.Context, trm.Settings, func(ctx context.Context) error) error
}

type ControllerOrder struct {
	storage  Storage
	changers map[model.TypePacking]order_changers.ChangerOrder
	tm       Manager
}

func New(repository Storage, changers map[model.TypePacking]order_changers.ChangerOrder, tm Manager) *ControllerOrder {
	return &ControllerOrder{
		storage:  repository,
		changers: changers,
		tm:       tm,
	}
}
