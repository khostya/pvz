package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/cache"
	"github.com/khostya/pvz/internal/domain"
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
		cache      cache.Cache[string, *pvz_v1.GetPVZListResponse]
	}

	pvzService interface {
		GetAllPvzList(ctx context.Context) ([]*domain.PVZ, error)
	}

	Deps struct {
		PvzService         pvzService
		GetPVZListResponse cache.Cache[string, *pvz_v1.GetPVZListResponse]
	}
)

var (
	cacheKey = uuid.New().String()
)

func NewServer(service Deps) *Server {
	return &Server{
		pvzService: service.PvzService,
		cache:      service.GetPVZListResponse,
	}
}

func (s *Server) Register(server grpc.ServiceRegistrar) {
	pvz_v1.RegisterPVZServiceServer(server, s)
}

func (s *Server) GetPVZList(ctx context.Context, _ *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	if v, ok := s.cache.Get(cacheKey); ok {
		return v, nil
	}

	pvzs, err := s.pvzService.GetAllPvzList(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pvzList := pvzListToResponse(pvzs)

	pvz := &pvz_v1.GetPVZListResponse{Pvzs: pvzList}
	s.cache.Put(cacheKey, pvz)
	return pvz, nil
}

func pvzListToResponse(pvzs []*domain.PVZ) []*pvz_v1.PVZ {
	pvzList := make([]*pvz_v1.PVZ, 0, len(pvzs))
	for _, pvz := range pvzs {
		pvzList = append(pvzList, &pvz_v1.PVZ{
			Id:               pvz.ID.String(),
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
			City:             string(pvz.City),
		})
	}

	return pvzList
}
