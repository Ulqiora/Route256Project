package pickpoint

import (
	"context"

	"homework/internal/database/cache"
	"homework/internal/model"
	"homework/internal/service/broker_io"
)

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Controller interface {
	Create(ctx context.Context, object model.PickPoint) (uint64, error)
	GetByID(ctx context.Context, id uint64) (model.PickPoint, error)
	List(ctx context.Context) ([]model.PickPoint, error)
	Update(ctx context.Context, object model.PickPoint) (uint64, error)
	Delete(ctx context.Context, id uint64) error
}

type Service struct {
	controller Controller
	sender     Sender
	cacher     cache.Cache
}

func New(controller Controller, sender Sender, cacher cache.Cache) *Service {
	service := Service{
		controller: controller,
		sender:     sender,
		cacher:     cacher,
	}
	return &service
}
