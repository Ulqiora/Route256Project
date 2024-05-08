package order_courier

import (
	"context"
	"net/http"

	"homework/internal/model"
	"homework/internal/service/broker_io"
)

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Controller interface {
	AcceptOrder(ctx context.Context, data model.OrderInitData) (uint64, error)
	ReturnOrderToCourier(ctx context.Context, idOrder uint64) error
}

// CourierService - service заказов, которые доставляет или забирает курьер
type CourierService interface {
	// AcceptOrder Курьер доставляет заказ на пункт выдачи
	AcceptOrder(w http.ResponseWriter, req *http.Request)
	// ReturnOrderToCourier Выдача заказа курьеру для возврата его на склад
	ReturnOrderToCourier(w http.ResponseWriter, req *http.Request)
}
type Service struct {
	controller Controller
	sender     Sender
}

func New(controller Controller, sender Sender) *Service {
	return &Service{controller: controller, sender: sender}
}
