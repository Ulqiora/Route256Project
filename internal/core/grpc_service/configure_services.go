package grpc_service

import (
	"context"
	"fmt"
	"log/slog"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"homework/internal/config"
	"homework/internal/database/cache/redis"
	"homework/internal/infrastructure/kafka"
	"homework/internal/service/broker_io"
	"homework/internal/service/grpc/order/order_client"
	"homework/internal/service/grpc/order/order_courier"
	"homework/internal/service/grpc/order/order_delivery"
	"homework/internal/service/grpc/pickpoint"
)

type Services struct {
	pickpointSrv     *pickpoint.Service
	orderClientSrv   *order_client.Service
	orderCourierSrv  *order_courier.Service
	orderDeliverySrv *order_delivery.Service
}

func ConfigureServices(controller Controllers, config *config.Config, metrics *grpc_prometheus.ServerMetrics) (*prometheus.Registry, *grpc.Server, error) {
	kafkaProducer, err := kafka.NewProducer(config.Kafka, slog.Default(), slog.Default())
	if err != nil {
		fmt.Println(err)
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
	return registry, grpcServer, nil
}
