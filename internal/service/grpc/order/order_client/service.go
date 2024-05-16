package order_client

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/api"
	"github.com/Ulqiora/Route256Project/internal/monitoring/prometheus_cli/metrics"
	"github.com/Ulqiora/Route256Project/internal/service/broker_io"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Controller interface {
	IssuingToCustomer(ctx context.Context, idOrders []string) error
	ReturnOrder(ctx context.Context, orderId string, customerId string) error
}

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Service struct {
	api.OrderClientServer
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
	api.RegisterOrderClientServer(server, New(controller, sender, tracer))
	registry.MustRegister(service.counterIssuedOrders)
}

func (s *Service) IssuingAnOrderCustomer(ctx context.Context, req *api.IssuingAnOrderCustomerRequest) (*api.IssuingAnOrderCustomerResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"IssuingAnOrderCustomer")
	defer span.End()

	if req.String() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty request. Set ids of orders for issuing")
	}
	err := s.controller.IssuingToCustomer(ctx, ConvertUUIDsToStrings(req.OrderIds))
	if err == nil {
		s.counterIssuedOrders.Add(float64(len(req.OrderIds)))
	}
	return &api.IssuingAnOrderCustomerResponse{}, err
}

func (s *Service) ReturnOrder(ctx context.Context, req *api.ReturnOrderRequest) (*api.ReturnOrderResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"ReturnOrder")
	defer span.End()
	if req.String() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty request. Set id of orders for returning")
	}
	err := s.controller.ReturnOrder(ctx, req.OrderId.Value, req.CustomerId.Value)
	return &api.ReturnOrderResponse{}, err
}

func ConvertUUIDsToStrings(ids []*api.UUID) []string {
	var result []string
	for i := range ids {
		result = append(result, ids[i].Value)
	}
	return result
}
