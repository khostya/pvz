package pvz

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	mock_postgres "github.com/khostya/pvz/internal/repo/postgres/mocks"
	mock_transactor "github.com/khostya/pvz/pkg/postgres/transactor/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

var (
	ErrPvzOops = errors.New("pvz oops error")
)

type productMocks struct {
	pvzRepo    *mock_postgres.MockPvzRepo
	transactor *mock_transactor.MockTransactor
}

func newMocks(t *testing.T) productMocks {
	ctrl := gomock.NewController(t)
	return productMocks{
		transactor: mock_transactor.NewMockTransactor(ctrl),
		pvzRepo:    mock_postgres.NewMockPvzRepo(ctrl),
	}
}

func newDepsUseCase(mocks productMocks) DepsUseCase {
	return DepsUseCase{
		TransactionManager: mocks.transactor,
		PvzRepo:            mocks.pvzRepo,
	}
}

func TestPvzUseCase_Create(t *testing.T) {
	t.Parallel()

	type test struct {
		name    string
		input   dto.CreatePvzParam
		pvz     *domain.PVZ
		mockFn  func(ctx context.Context, test test, m productMocks)
		wantErr error
	}

	input := dto.CreatePvzParam{
		ID:               uuid.New(),
		City:             domain.CityMoscow,
		RegistrationDate: time.Now(),
		CreatorRole:      domain.UserRoleModerator,
	}

	ctx := context.Background()
	tests := []test{
		{
			name:    "ok",
			input:   input,
			wantErr: nil,
			pvz:     &domain.PVZ{ID: uuid.New()},
			mockFn: func(ctx context.Context, test test, m productMocks) {
				m.pvzRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.pvz, nil)
			},
		},
		{
			name:    "error moderator only",
			input:   dto.CreatePvzParam{CreatorRole: domain.UserRoleEmployee},
			wantErr: domain.ErrModeratorOnly,
			mockFn: func(ctx context.Context, test test, m productMocks) {
			},
		},
		{
			name:    "error create pvz",
			input:   input,
			wantErr: ErrPvzOops,
			pvz:     nil,
			mockFn: func(ctx context.Context, test test, m productMocks) {
				m.pvzRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, ErrPvzOops)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)

			useCase := NewUseCase(newDepsUseCase(mocks))

			tt.mockFn(ctx, tt, mocks)

			pvz, err := useCase.Create(ctx, tt.input)
			require.Equal(t, tt.wantErr, err)

			require.Equal(t, tt.pvz, pvz)
		})
	}
}

func runRepeatableReadAndUnwrap(m productMocks) {
	runRepeatableRead(m)
	m.transactor.EXPECT().Unwrap(gomock.Any()).
		Times(1).
		DoAndReturn(func(err error) error {
			return err
		})
}

func runRepeatableRead(m productMocks) {
	m.transactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).
		Times(1).
		DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
			return transaction(ctx)
		})
}
