package grpc

import (
	"errors"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	mock_pvz "github.com/khostya/pvz/internal/usecase/pvz/mocks"
	pvz_v1 "github.com/khostya/pvz/pkg/api/v1/proto"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

type mocks struct {
	pvz *mock_pvz.MockPvz
}

func newMocks(t *testing.T) mocks {
	ctrl := gomock.NewController(t)
	return mocks{
		pvz: mock_pvz.NewMockPvz(ctrl),
	}
}

func NewMockServer(m mocks) *Server {
	return NewServer(m.pvz)
}

var (
	errOops = errors.New("oops")
)

func TestAuth_PostDummyLogin(t *testing.T) {
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
				m.pvz.EXPECT().GetAllPvzs(gomock.Any()).
					Times(1).
					Return(test.pvzList, nil)
			},
		},
		{
			name: "dummy login error",
			err:  status.Error(codes.Internal, errOops.Error()),
			mockFn: func(test *test, m mocks) {
				m.pvz.EXPECT().GetAllPvzs(gomock.Any()).
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
