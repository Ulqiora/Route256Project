package order_delivery

import (
	"context"
	"net/http"

	"homework/internal/controller"
	"homework/internal/database/cache"
	"homework/internal/model"
	"homework/internal/service/broker_io"
)

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Controller interface {
	SearchOrders(ctx context.Context, customerID uint64, values controller.ValuesView) ([]model.Order, error)
	GetReturnedOrders(ctx context.Context, values controller.ValuesView) ([]model.Order, error)
}

// DeliveryService Сервис, который больше работает с ПВЗ
type DeliveryService interface {
	// GetReturnedOrders получить список заказов с учетом пагинации
	GetReturnedOrders(w http.ResponseWriter, req *http.Request)
	// SearchOrders Поиск заказов пользователя с фильтрацией
	SearchOrders(w http.ResponseWriter, req *http.Request)
}

type Service struct {
	controller Controller
	sender     Sender
	cacher     cache.Cache
}

func New(storage Controller, sender Sender, cacher cache.Cache) *Service {
	return &Service{
		controller: storage,
		sender:     sender,
		cacher:     cacher,
	}
}
