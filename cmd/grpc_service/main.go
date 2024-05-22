package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	core "github.com/Ulqiora/Route256Project/internal/core/grpc_service"
	"github.com/Ulqiora/Route256Project/internal/monitoring/jaeger"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		ctx = context.Background()
	)
	configApp := core.LoadConfig(os.Args[1])
	tracerFunc, err := jaeger.ConfigureTracer(ctx, configApp.Monitoring)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := tracerFunc(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()
	manager, bd, err := core.ConfigureTransactionManager(ctx, configApp)
	if err != nil {
		fmt.Println(err)
		return
	}
	repositories, err := core.ConfigureRepositories(ctx, bd)
	if err != nil {
		fmt.Println(err)
		return
	}
	controllers := core.ConfigureControllers(repositories, manager)
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	registry, grpcServer, err := core.ConfigureServices(controllers, configApp, grpcMetrics)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		err = http.ListenAndServe(configApp.Monitoring.Prometheus.Address, promhttp.HandlerFor(registry, promhttp.HandlerOpts{EnableOpenMetrics: true}))
		fmt.Println(err)
	}()

	list, err := net.Listen("tcp", configApp.Grpc.Address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer list.Close()
	if err = grpcServer.Serve(list); err != nil {
		fmt.Println(err)
		return
	}
}
