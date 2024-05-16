package order_delivery

import (
	"context"
	"errors"
	"strconv"

	"github.com/Ulqiora/Route256Project/internal/api"
	"github.com/Ulqiora/Route256Project/internal/controller"
	"github.com/Ulqiora/Route256Project/internal/database/cache"
	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/service/broker_io"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Controller interface {
	SearchOrders(ctx context.Context, customerID string, values controller.ValuesView) ([]model.Order, error)
	GetReturnedOrders(ctx context.Context, values controller.ValuesView) ([]model.Order, error)
}

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Service struct {
	api.OrderDeliveryServer
	controller Controller
	sender     Sender
	cacher     cache.Cache
	tracer     trace.Tracer
}

func New(controller Controller, sender Sender, cacher cache.Cache, tracer trace.Tracer) *Service {
	return &Service{
		controller: controller,
		sender:     sender,
		cacher:     cacher,
		tracer:     tracer,
	}
}

func RegisterService(server *grpc.Server, controller Controller, sender Sender, cacher cache.Cache, tracer trace.Tracer) {
	api.RegisterOrderDeliveryServer(server, New(controller, sender, cacher, tracer))
}

func (s *Service) SearchOrders(ctx context.Context, req *api.SearchOrdersRequest) (*api.SearchOrdersResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"SearchOrders")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("request param is empty")
	}
	view := &controller.ParametersView{
		Data: map[string]string{
			"last_n":       strconv.Itoa(int(req.LastN)),
			"pickpoint_id": req.PickpointId.Value,
		},
	}
	orders, err := s.controller.SearchOrders(ctx, req.CustomerId.Value, view)
	if err != nil {
		return nil, err
	}
	response := &api.SearchOrdersResponse{}
	for _, order := range orders {
		response.Orders = append(response.Orders, &api.Order{
			ID:           &api.UUID{Value: order.ID},
			Customer_ID:  &api.UUID{Value: order.CustomerID},
			Pickpoint_ID: &api.UUID{Value: order.PickPointID},
			ShelfTime:    timestamppb.New(*order.ShelfLife.Time()),
			TimeCreated:  timestamppb.New(*order.TimeCreated.Time()),
			DateReceipt:  timestamppb.New(*order.DateReceipt.Time()),
			Penny:        order.Penny,
			Weight:       order.Weight,
			State:        order.State,
		})
	}
	return response, nil
}

func (s *Service) GetReturnedOrders(ctx context.Context, req *api.GetReturnedOrdersRequest) (*api.GetReturnedOrdersResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"GetReturnedOrders")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("request param is empty")
	}
	view := &controller.ParametersView{
		Data: map[string]string{
			"page":  strconv.Itoa(int(req.Page)),
			"limit": strconv.Itoa(int(req.Limit)),
		},
	}
	orders, err := s.controller.GetReturnedOrders(ctx, view)
	if err != nil {
		return nil, err
	}
	response := &api.GetReturnedOrdersResponse{}
	for _, order := range orders {
		response.Orders = append(response.Orders, &api.Order{
			ID:           &api.UUID{Value: order.ID},
			Customer_ID:  &api.UUID{Value: order.CustomerID},
			Pickpoint_ID: &api.UUID{Value: order.PickPointID},
			ShelfTime:    timestamppb.New(*order.ShelfLife.Time()),
			TimeCreated:  timestamppb.New(*order.TimeCreated.Time()),
			DateReceipt:  timestamppb.New(*order.DateReceipt.Time()),
			Penny:        order.Penny,
			Weight:       order.Weight,
			State:        order.State,
		})
	}
	return response, nil
}
