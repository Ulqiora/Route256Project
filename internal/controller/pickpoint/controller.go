package pickpoint

import (
	"context"

	"github.com/jackc/pgx/v4"
	"homework/internal/repository"
)

type Storage interface {
	GetByID(ctx context.Context, id int) (repository.PickPointDTO, error)
	List(ctx context.Context) ([]repository.PickPointDTO, error)
	Create(ctx context.Context, dto repository.PickPointDTO) (uint64, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, dto repository.PickPointDTO) (int, error)
}

type TransactionManager interface {
	Run(ctx context.Context, options pgx.TxOptions, f func(ctxTX context.Context) error) error
}

type Controller struct {
	storage Storage
	tm      TransactionManager
}

func New(storage Storage, tm TransactionManager) *Controller {
	return &Controller{
		storage: storage,
		tm:      tm,
	}
}
