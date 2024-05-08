package order_client

import (
	"context"
	"net/http"

	"homework/internal/service/broker_io"
)

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Controller interface {
	IssuingToCustomer(ctx context.Context, idOrders []uint64) error
	ReturnOrder(ctx context.Context, idOrder uint64, customer_id uint64) error
}

// ClientService - service заказов, которые возвращает или забирает клиент
type ClientService interface {
	// IssuingAnOrderCustomer выдача заказов клиенту
	IssuingAnOrderCustomer(w http.ResponseWriter, req *http.Request)
	// ReturnOrder Принять возврат заказа от клиента
	ReturnOrder(w http.ResponseWriter, req *http.Request)
}

type Service struct {
	controller Controller
	sender     Sender
}

func New(controller Controller, sender Sender) *Service {
	return &Service{
		sender:     sender,
		controller: controller,
	}
}
