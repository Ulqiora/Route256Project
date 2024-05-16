package client

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/jackc/pgtype"
)

type Storage interface {
	Create(ctx context.Context, dto repository.ClientDTO) (pgtype.UUID, error)
	GetByID(ctx context.Context, id pgtype.UUID) (repository.ClientDTO, error)
	Update(ctx context.Context, dto repository.ClientDTO) (pgtype.UUID, error)
	List(ctx context.Context) ([]repository.ClientDTO, error)
	Delete(ctx context.Context, client pgtype.UUID) error
}

type Manager interface {
	// Do processes a transaction inside a closure.
	Do(context.Context, func(ctx context.Context) error) error
	// DoWithSettings processes a transaction inside a closure with custom trm.Settings.
	DoWithSettings(context.Context, trm.Settings, func(ctx context.Context) error) error
}

type Controller struct {
	storage Storage
	tm      Manager
}

func New(storage Storage, tm Manager) *Controller {
	return &Controller{
		storage: storage,
		tm:      tm,
	}
}
