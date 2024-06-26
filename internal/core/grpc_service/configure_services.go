package grpc_service

import (
	"context"
	"log/slog"

	"github.com/Ulqiora/Route256Project/internal/config"
	"github.com/Ulqiora/Route256Project/internal/database/cache/redis"
	"github.com/Ulqiora/Route256Project/internal/infrastructure/kafka"
	"github.com/Ulqiora/Route256Project/internal/service/broker_io"
	"github.com/Ulqiora/Route256Project/internal/service/grpc/client"
	"github.com/Ulqiora/Route256Project/internal/service/grpc/order/order_client"
	"github.com/Ulqiora/Route256Project/internal/service/grpc/order/order_courier"
	"github.com/Ulqiora/Route256Project/internal/service/grpc/order/order_delivery"
	"github.com/Ulqiora/Route256Project/internal/service/grpc/pickpoint"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Services struct {
	pickpointSrv     *pickpoint.Service
	orderClientSrv   *order_client.Service
	orderCourierSrv  *order_courier.Service
	orderDeliverySrv *order_delivery.Service
	clientSrv        *client.Service
}

func ConfigureServices(controller Controllers, config *config.Config, metrics *grpc_prometheus.ServerMetrics) (*prometheus.Registry, *grpc.Server, error) {
	kafkaProducer, err := kafka.NewProducer(config.Kafka, slog.Default(), slog.Default())
	if err != nil {
		return nil, nil, err
	}
	kafkaProducer.Run()
	redisClient := redis.New(config.Redis)
	if err = redisClient.Ping(context.Background()); err != nil {
		return nil, nil, err
	}
	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			metrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			metrics.StreamServerInterceptor(),
		),
	)
	registry := prometheus.NewRegistry()
	reflection.Register(grpcServer)
	metrics.InitializeMetrics(grpcServer)
	registry.MustRegister(metrics)
	order_client.RegisterService(grpcServer, controller.orderCtrl, broker_io.NewKafkaSender(kafkaProducer, "order"), otel.Tracer("order-client_repo-service"), registry)
	order_delivery.RegisterService(grpcServer, controller.orderCtrl, broker_io.NewKafkaSender(kafkaProducer, "order"), redisClient, otel.Tracer("order-delivery-service"))
	order_courier.RegisterService(grpcServer, controller.orderCtrl, broker_io.NewKafkaSender(kafkaProducer, "order"), otel.Tracer("order-courier-service"))
	pickpoint.RegisterService(grpcServer, controller.pickpointCtrl, broker_io.NewKafkaSender(kafkaProducer, "pickpoint"), redisClient, otel.Tracer("pickpoint-service"))
	client.RegisterService(grpcServer, controller.clientCtrl, broker_io.NewKafkaSender(kafkaProducer, "client"), redisClient, otel.Tracer("client-service"))
	return registry, grpcServer, nil
}
