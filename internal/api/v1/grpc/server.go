package grpc

import (
	"context"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	pvz_v1 "github.com/khostya/pvz/pkg/api/v1/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	Server struct {
		pvz_v1.UnimplementedPVZServiceServer

		pvzService pvzService
	}

	pvzService interface {
		GetPvz(ctx context.Context, pvz dto.GetPvzParam) ([]*domain.PVZ, error)
	}
)

func NewServer(service pvzService) *Server {
	return &Server{
		pvzService: service,
	}
}

func (s *Server) Register(server grpc.ServiceRegistrar) {
	pvz_v1.RegisterPVZServiceServer(server, s)
}
func (s *Server) GetPVZList(ctx context.Context, _ *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	pvzs, err := s.pvzService.GetPvz(ctx, dto.GetPvzParam{})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := make([]*pvz_v1.PVZ, 0, len(pvzs))
	for _, pvz := range pvzs {
		resp = append(resp, &pvz_v1.PVZ{
			Id:               pvz.ID.String(),
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
			City:             string(pvz.City),
		})
	}

	return &pvz_v1.GetPVZListResponse{Pvzs: resp}, nil
}
