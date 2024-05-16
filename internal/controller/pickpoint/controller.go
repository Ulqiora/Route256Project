package pickpoint

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/repository"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/jackc/pgtype"
)

type Storage interface {
	GetByID(ctx context.Context, id pgtype.UUID) (repository.PickPointDTO, error)
	List(ctx context.Context) ([]repository.PickPointDTO, error)
	Create(ctx context.Context, dto repository.PickPointDTO) (pgtype.UUID, error)
	Delete(ctx context.Context, id pgtype.UUID) error
	Update(ctx context.Context, dto repository.PickPointDTO) (pgtype.UUID, error)
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
