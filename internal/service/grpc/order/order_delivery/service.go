package order_delivery

import (
	"context"
	"errors"
	"strconv"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework/internal/controller"
	"homework/internal/database/cache"
	pb "homework/internal/gen_proto"
	"homework/internal/model"
	"homework/internal/service/broker_io"
)

type Controller interface {
	SearchOrders(ctx context.Context, customerID uint64, values controller.ValuesView) ([]model.Order, error)
	GetReturnedOrders(ctx context.Context, values controller.ValuesView) ([]model.Order, error)
}

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Service struct {
	pb.OrderDeliveryServer
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
	pb.RegisterOrderDeliveryServer(server, New(controller, sender, cacher, tracer))
}

func (s *Service) SearchOrders(ctx context.Context, req *pb.SearchOrdersRequest) (*pb.SearchOrdersResponse, error) {
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
			"pickpoint_id": strconv.Itoa(int(req.PickpointId)),
		},
	}
	orders, err := s.controller.SearchOrders(ctx, uint64(req.CustomerId), view)
	if err != nil {
		return nil, err
	}
	response := &pb.SearchOrdersResponse{}
	for _, order := range orders {
		response.Orders = append(response.Orders, &pb.Order{
			ID:           order.ID,
			Customer_ID:  order.CustomerID,
			Pickpoint_ID: order.PickPointID,
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

func (s *Service) GetReturnedOrders(ctx context.Context, req *pb.GetReturnedOrdersRequest) (*pb.GetReturnedOrdersResponse, error) {
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
	response := &pb.GetReturnedOrdersResponse{}
	for _, order := range orders {
		response.Orders = append(response.Orders, &pb.Order{
			ID:           order.ID,
			Customer_ID:  order.CustomerID,
			Pickpoint_ID: order.PickPointID,
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
