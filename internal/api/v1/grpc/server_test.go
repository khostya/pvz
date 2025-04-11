package grpc

import (
	"errors"
	"github.com/google/uuid"
	mock_cache "github.com/khostya/pvz/internal/cache/mocks"
	"github.com/khostya/pvz/internal/domain"
	mock_pvz "github.com/khostya/pvz/internal/usecase/pvz/mocks"
	pvz_v1 "github.com/khostya/pvz/pkg/api/v1/proto"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

type mocks struct {
	pvz   *mock_pvz.MockPvz
	cache *mock_cache.MockCache[string, *pvz_v1.GetPVZListResponse]
}

func newMocks(t *testing.T) mocks {
	ctrl := gomock.NewController(t)
	return mocks{
		pvz:   mock_pvz.NewMockPvz(ctrl),
		cache: mock_cache.NewMockCache[string, *pvz_v1.GetPVZListResponse](ctrl),
	}
}

func NewMockServer(m mocks) *Server {
	return NewServer(Deps{
		PvzService:         m.pvz,
		GetPVZListResponse: m.cache,
	})
}

var (
	errOops = errors.New("oops")
)

func TestServer_GetPVZList(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	type test struct {
		name    string
		input   *pvz_v1.GetPVZListRequest
		err     error
		pvzList []*domain.PVZ
		mockFn  func(test *test, m mocks)
	}

	pvz := &domain.PVZ{City: domain.CityMoscow, ID: uuid.New(), RegistrationDate: time.Now()}

	tests := []*test{
		{
			name:    "ok",
			pvzList: []*domain.PVZ{pvz},
			mockFn: func(test *test, m mocks) {
				m.cache.EXPECT().Get(gomock.Any()).
					Times(1).
					Return(nil, false)
				m.pvz.EXPECT().GetAllPvzList(gomock.Any()).
					Times(1).
					Return(test.pvzList, nil)
				m.cache.EXPECT().Put(gomock.Any(), gomock.Any()).
					Times(1).
					Return()
			},
		},
		{
			name:    "cache hit",
			pvzList: []*domain.PVZ{pvz},
			mockFn: func(test *test, m mocks) {
				m.cache.EXPECT().Get(gomock.Any()).
					Times(1).
					Return(&pvz_v1.GetPVZListResponse{Pvzs: pvzListToResponse(test.pvzList)}, true)
			},
		},
		{
			name: "dummy login error",
			err:  status.Error(codes.Internal, errOops.Error()),
			mockFn: func(test *test, m mocks) {
				m.cache.EXPECT().Get(gomock.Any()).
					Times(1).
					Return(nil, false)
				m.pvz.EXPECT().GetAllPvzList(gomock.Any()).
					Times(1).
					Return(nil, errOops)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)

			tt.mockFn(tt, mocks)

			server := NewMockServer(mocks)

			res, err := server.GetPVZList(ctx, tt.input)
			require.Equal(t, tt.err, err)
			if err != nil {
				return
			}

			pvzList := make([]*domain.PVZ, 0)
			for _, item := range res.GetPvzs() {
				pvzList = append(pvzList, &domain.PVZ{
					ID:               uuid.MustParse(item.GetId()),
					RegistrationDate: item.GetRegistrationDate().AsTime(),
					City:             domain.City(item.GetCity()),
				})
			}

			require.EqualExportedValues(t, tt.pvzList, pvzList)
		})
	}
}

func TestServer_Register(t *testing.T) {
	t.Parallel()

	mocks := newMocks(t)
	server := NewMockServer(mocks)

	server.Register(grpc.NewServer())
}
