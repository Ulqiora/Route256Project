package order_courier

import (
	"context"
	"errors"

	"github.com/Ulqiora/Route256Project/internal/api"
	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/service/broker_io"
	jtime "github.com/Ulqiora/Route256Project/pkg/wrapper/jsontime"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type Controller interface {
	AcceptOrder(ctx context.Context, data model.OrderInitData) (string, error)
	ReturnOrderToCourier(ctx context.Context, orderID string) error
}
type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Service struct {
	api.OrderCourierServer
	controller Controller
	sender     Sender
	tracer     trace.Tracer
}

func New(controller Controller, sender Sender, tracer trace.Tracer) *Service {
	return &Service{
		sender:     sender,
		controller: controller,
		tracer:     tracer,
	}
}

func RegisterService(server *grpc.Server, controller Controller, sender Sender, tracer trace.Tracer) {
	api.RegisterOrderCourierServer(server, New(controller, sender, tracer))
}

func (s *Service) AcceptOrder(ctx context.Context, req *api.AcceptOrderRequest) (*api.AcceptOrderResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"AcceptOrder")
	defer span.End()
	if req.Order.String() == "" {
		return nil, errors.New("empty order")
	}
	orderID, err := s.controller.AcceptOrder(ctx, model.OrderInitData{
		CustomerID:  req.Order.Customer_ID.Value,
		PickPointID: req.Order.Pickpoint_ID.Value,
		ShelfLife:   jtime.TimeWrap(req.Order.ShelfTime.AsTime()),
		Penny:       req.Order.Penny,
		Weight:      req.Order.Weight,
		Type:        model.TypePacking(req.Order.TypePacking),
	})
	return &api.AcceptOrderResponse{OrderId: &api.UUID{
		Value: orderID,
	}}, err
}
func (s *Service) ReturnOrderToCourier(ctx context.Context, req *api.ReturnOrderToCourierRequest) (*api.ReturnOrderToCourierResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"ReturnOrderToCourier")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("request is empty")
	}
	err := s.controller.ReturnOrderToCourier(ctx, req.OrderId.Value)
	return &api.ReturnOrderToCourierResponse{}, err
}
