package order_courier

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	pb "homework/internal/gen_proto"
	"homework/internal/model"
	"homework/internal/service/broker_io"
	jtime "homework/pkg/wrapper/jsontime"
)

type Controller interface {
	AcceptOrder(ctx context.Context, data model.OrderInitData) (uint64, error)
	ReturnOrderToCourier(ctx context.Context, idOrder uint64) error
}
type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Service struct {
	pb.OrderCourierServer
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
	pb.RegisterOrderCourierServer(server, New(controller, sender, tracer))
}

func (s *Service) AcceptOrder(ctx context.Context, req *pb.AcceptOrderRequest) (*pb.AcceptOrderResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"AcceptOrder")
	defer span.End()
	if req.Order.String() == "" {
		return nil, errors.New("empty order")
	}
	orderID, err := s.controller.AcceptOrder(ctx, model.OrderInitData{
		CustomerID:  req.Order.Customer_ID,
		PickPointID: req.Order.Pickpoint_ID,
		ShelfLife:   jtime.TimeWrap(req.Order.ShelfTime.AsTime()),
		Penny:       req.Order.Penny,
		Weight:      req.Order.Weight,
		Type:        model.TypePacking(req.Order.TypePacking),
	})
	return &pb.AcceptOrderResponse{OrderId: int64(orderID)}, err
}
func (s *Service) ReturnOrderToCourier(ctx context.Context, req *pb.ReturnOrderToCourierRequest) (*pb.ReturnOrderToCourierResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"ReturnOrderToCourier")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("request is empty")
	}
	err := s.controller.ReturnOrderToCourier(ctx, uint64(req.OrderId))
	return &pb.ReturnOrderToCourierResponse{}, err
}
