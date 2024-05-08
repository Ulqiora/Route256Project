package pickpoint

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"homework/internal/database/cache"
	pb "homework/internal/gen_proto"
	"homework/internal/model"
	"homework/internal/service/broker_io"
)

type Controller interface {
	Create(ctx context.Context, object model.PickPoint) (uint64, error)
	GetByID(ctx context.Context, id uint64) (model.PickPoint, error)
	List(ctx context.Context) ([]model.PickPoint, error)
	Update(ctx context.Context, object model.PickPoint) (uint64, error)
	Delete(ctx context.Context, id uint64) error
}

type Sender interface {
	SendMessage(message broker_io.RequestMessage)
}

type Service struct {
	pb.UnimplementedPickPointServiceServer
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
	pb.RegisterPickPointServiceServer(server, New(controller, sender, cacher, tracer))
}

func (s *Service) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"CreatePickpoint")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("empty request")
	}
	pickpoint := model.PickPoint{
		ID:             int(req.Pickpoint.ID),
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
	return &pb.CreateResponse{Pickpoint_ID: int32(id)}, nil
}
func (s *Service) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"GetPickpoint")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("empty request")
	}
	pickpoint, err := s.controller.GetByID(ctx, uint64(req.Pickpoint_ID))
	if err != nil {
		return nil, err
	}
	grpcPickpoint := GetGrpcPickpoint(pickpoint)
	return &pb.GetResponse{Pickpoint: grpcPickpoint}, nil
}
func (s *Service) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"ListPickpoint")
	defer span.End()
	pickpoints, err := s.controller.List(ctx)
	if err != nil {
		return nil, err
	}
	response := &pb.ListResponse{}
	for _, pp := range pickpoints {
		response.Pickpoint = append(response.Pickpoint, GetGrpcPickpoint(pp))
	}
	return response, nil
}
func (s *Service) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"UpdatePickpoint")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("empty request")
	}
	pickpoint := model.PickPoint{
		ID:             int(req.Pickpoint.ID),
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
	return &pb.UpdateResponse{
		Pickpoint_ID: int32(id),
	}, err
}
func (s *Service) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	ctx, span := s.tracer.Start(
		ctx,
		"DeletePickpoint")
	defer span.End()
	if req.String() == "" {
		return nil, errors.New("empty request")
	}
	err := s.controller.Delete(ctx, uint64(req.Pickpoint_ID))
	return &pb.DeleteResponse{}, err
}

func GetGrpcPickpoint(pickpoint model.PickPoint) *pb.PickPoint {
	result := &pb.PickPoint{
		ID:      int32(pickpoint.ID),
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
