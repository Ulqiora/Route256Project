package order_client

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "homework/internal/gen_proto"
	"homework/internal/monitoring/prometheus_cli/metrics"
	"homework/internal/service/broker_io"
)

type Controller interface {
	IssuingToCustomer(ctx context.Context, idOrders []uint64) error
	ReturnOrder(ctx context.Context, orderId uint64, customerId uint64) error
}

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Service struct {
	pb.OrderClientServer
	controller          Controller
	sender              Sender
	tracer              trace.Tracer
	counterIssuedOrders prometheus.Counter
}

func New(controller Controller, sender Sender, tracer trace.Tracer) *Service {
	return &Service{
		sender:              sender,
		controller:          controller,
		tracer:              tracer,
		counterIssuedOrders: metrics.NewCounterIssuedOrders(),
	}
}

func RegisterService(server *grpc.Server, controller Controller, sender Sender, tracer trace.Tracer, registry *prometheus.Registry) {
	service := New(controller, sender, tracer)
	pb.RegisterOrderClientServer(server, New(controller, sender, tracer))
	registry.MustRegister(service.counterIssuedOrders)
}

func (s *Service) IssuingAnOrderCustomer(ctx context.Context, req *pb.IssuingAnOrderCustomerRequest) (*pb.IssuingAnOrderCustomerResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"IssuingAnOrderCustomer")
	defer span.End()

	if req.String() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty request. Set ids of orders for issuing")
	}
	err := s.controller.IssuingToCustomer(ctx, req.OrderIds)
	if err == nil {
		s.counterIssuedOrders.Add(float64(len(req.OrderIds)))
	}
	return &pb.IssuingAnOrderCustomerResponse{}, err
}

func (s *Service) ReturnOrder(ctx context.Context, req *pb.ReturnOrderRequest) (*pb.ReturnOrderResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"ReturnOrder")
	defer span.End()
	if req.String() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty request. Set id of orders for returning")
	}
	err := s.controller.ReturnOrder(ctx, req.OrderId, req.CustomerId)
	return &pb.ReturnOrderResponse{}, err
}
