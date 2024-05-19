package pickpoint

import (
	"context"
	"errors"

	"github.com/Ulqiora/Route256Project/internal/api"
	pb "github.com/Ulqiora/Route256Project/internal/api"
	"github.com/Ulqiora/Route256Project/internal/database/cache"
	"github.com/Ulqiora/Route256Project/internal/model"
	"github.com/Ulqiora/Route256Project/internal/service/broker_io"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type Controller interface {
	Create(ctx context.Context, object model.PickPoint) (string, error)
	GetByID(ctx context.Context, id string) (model.PickPoint, error)
	List(ctx context.Context) ([]model.PickPoint, error)
	Update(ctx context.Context, object model.PickPoint) (string, error)
	Delete(ctx context.Context, id string) error
}

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Service struct {
	api.UnimplementedPickPointServiceServer
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
	api.RegisterPickPointServiceServer(server, New(controller, sender, cacher, tracer))
}

func (s *Service) Create(ctx context.Context, req *api.CreatePickPointRequest) (*api.CreatePickPointResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"CreatePickpoint")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("empty request")
	}
	pickpoint := model.PickPoint{
		Name:           req.Pickpoint.Name,
		Address:        req.Pickpoint.Address,
		ContactDetails: nil,
	}
	for _, contact := range req.Pickpoint.ContactDetails {
		pickpoint.ContactDetails = append(pickpoint.ContactDetails, model.ContactDetail{
			Type:   contact.Type,
			Detail: contact.Detail,
		})
	}
	id, err := s.controller.Create(ctx, pickpoint)
	if err != nil {
		return nil, err
	}
	var response api.CreatePickPointResponse
	response.Pickpoint_ID = &api.UUID{
		Value: id,
	}
	return &response, nil
}

func (s *Service) Get(ctx context.Context, req *api.GetPickPointRequest) (*api.GetPickPointResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"GetPickpoint")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("empty request")
	}
	pickpoint, err := s.controller.GetByID(ctx, req.Pickpoint_ID.Value)
	if err != nil {
		return nil, err
	}
	grpcPickpoint := GetGrpcPickpoint(pickpoint)
	var response api.GetPickPointResponse
	response.Pickpoint = grpcPickpoint
	return &response, nil
}
func (s *Service) List(ctx context.Context, _ *api.ListPickPointRequest) (*api.ListPickPointResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"ListPickpoint")
	defer span.End()
	pickpoints, err := s.controller.List(ctx)
	if err != nil {
		return nil, err
	}
	response := &api.ListPickPointResponse{}
	for _, pp := range pickpoints {
		response.Pickpoint = append(response.Pickpoint, GetGrpcPickpoint(pp))
	}
	return response, nil
}
func (s *Service) Update(ctx context.Context, req *api.UpdatePickPointRequest) (*api.UpdatePickPointResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"UpdatePickpoint")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("empty request")
	}
	pickpoint := model.PickPoint{
		ID:             req.Pickpoint.ID.Value,
		Name:           req.Pickpoint.Name,
		Address:        req.Pickpoint.Address,
		ContactDetails: nil,
	}
	for _, contact := range req.Pickpoint.ContactDetails {
		pickpoint.ContactDetails = append(pickpoint.ContactDetails, model.ContactDetail{
			Type:   contact.Type,
			Detail: contact.Detail,
		})
	}
	id, err := s.controller.Update(ctx, pickpoint)
	var response api.UpdatePickPointResponse
	response.Pickpoint_ID.Value = id
	return &response, err
}
func (s *Service) Delete(ctx context.Context, req *api.DeletePickPointRequest) (*api.DeletePickPointResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"DeletePickpoint")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("empty request")
	}
	err := s.controller.Delete(ctx, req.Pickpoint_ID.Value)
	return &api.DeletePickPointResponse{}, err
}

func GetGrpcPickpoint(pickpoint model.PickPoint) *api.PickPoint {
	result := &pb.PickPoint{
		ID:      &pb.UUID{Value: pickpoint.ID},
		Name:    pickpoint.Name,
		Address: pickpoint.Address,
	}
	for _, details := range pickpoint.ContactDetails {
		result.ContactDetails = append(result.ContactDetails, &pb.ContactDetails{
			Type:   details.Type,
			Detail: details.Detail,
		})
	}
	return result
}
