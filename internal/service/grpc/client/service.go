package client

import (
	"context"

	"github.com/Ulqiora/Route256Project/internal/api"
	"github.com/Ulqiora/Route256Project/internal/database/cache"
	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/service/broker_io"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type Controller interface {
	Create(ctx context.Context, obj model.Client) (string, error)
	Delete(ctx context.Context, ObjId string) error
	Get(ctx context.Context, ObjId string) (model.Client, error)
	Update(ctx context.Context, obj model.Client) (string, error)
	List(ctx context.Context) ([]model.Client, error)
}

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Service struct {
	api.UnimplementedClientServiceServer
	controller Controller
	sender     Sender
	cacher     cache.Cache
	tracer     trace.Tracer
}

func New(controller Controller, sender Sender, cacher cache.Cache, tracer trace.Tracer) *Service {
	service := Service{
		controller: controller,
		sender:     sender,
		cacher:     cacher,
		tracer:     tracer,
	}
	return &service
}

func RegisterService(server *grpc.Server, controller Controller, sender Sender, cacher cache.Cache, tracer trace.Tracer) {
	api.RegisterClientServiceServer(server, New(controller, sender, cacher, tracer))
}

func (s *Service) Create(ctx context.Context, req *api.CreateClientRequest) (*api.CreateClientResponse, error) {
	var object model.Client
	err := object.LoadFromGrpcModel(req.Client)
	if err != nil {
		return nil, err
	}
	id, err := s.controller.Create(ctx, model.Client{
		Name: req.Client.Name,
	})
	var resp api.CreateClientResponse
	if err != nil {
		return &resp, err
	}
	resp.Client_ID = &api.UUID{
		Value: id,
	}
	return &resp, nil
}

func (s *Service) Get(ctx context.Context, req *api.GetClientRequest) (*api.GetClientResponse, error) {
	client, err := s.controller.Get(ctx, req.Client_ID.Value)
	if err != nil {
		return nil, err
	}
	var res api.GetClientResponse
	res.Client = client.MapToGrpcModel()
	return &res, nil

}
func (s *Service) List(ctx context.Context, req *api.ListClientRequest) (*api.ListClientResponse, error) {
	clients, err := s.controller.List(ctx)
	if err != nil {
		return nil, err
	}
	var res api.ListClientResponse
	for i := range clients {
		res.Client = append(res.Client, clients[i].MapToGrpcModel())
	}
	return &res, nil
}
func (s *Service) Update(ctx context.Context, req *api.UpdateClientRequest) (*api.UpdateClientResponse, error) {
	id, err := s.controller.Update(ctx, model.Client{
		ID:   req.Client.ID.Value,
		Name: req.Client.Name,
	})
	if err != nil {
		return nil, err
	}
	var resp api.UpdateClientResponse
	resp.Client_ID = &api.UUID{
		Value: id,
	}
	return &resp, nil
}
func (s *Service) Delete(ctx context.Context, res *api.DeleteClientRequest) (*api.DeleteClientResponse, error) {
	err := s.controller.Delete(ctx, res.Client_ID.Value)
	if err != nil {
		return nil, err
	}
	return &api.DeleteClientResponse{}, nil
}
