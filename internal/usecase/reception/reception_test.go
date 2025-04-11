package reception

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	mock_postgres "github.com/khostya/pvz/internal/repo/postgres/mocks"
	"github.com/khostya/pvz/pkg/postgres/repoerr"
	mock_transactor "github.com/khostya/pvz/pkg/postgres/transactor/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

var (
	ErrReceptionOops = errors.New("reception oops error")
)

type receptionMocks struct {
	receptionRepo *mock_postgres.MockReceptionRepo
	productRepo   *mock_postgres.MockProductRepo
	transactor    *mock_transactor.MockTransactor
}

func newMocks(t *testing.T) receptionMocks {
	ctrl := gomock.NewController(t)
	return receptionMocks{
		transactor:    mock_transactor.NewMockTransactor(ctrl),
		receptionRepo: mock_postgres.NewMockReceptionRepo(ctrl),
		productRepo:   mock_postgres.NewMockProductRepo(ctrl),
	}
}

func newDepsUseCase(mocks receptionMocks) DepsUseCase {
	return DepsUseCase{
		TransactionManager: mocks.transactor,
		ReceptionRepo:      mocks.receptionRepo,
		ProductRepo:        mocks.productRepo,
	}
}

func TestReceptionUseCase_Create(t *testing.T) {
	t.Parallel()

	type test struct {
		name      string
		input     dto.CreateReceptionParam
		reception *domain.Reception
		mockFn    func(test test, m receptionMocks)
		wantErr   error
	}

	input := dto.CreateReceptionParam{
		CreatorRole: domain.UserRoleEmployee,
		PvzID:       uuid.New(),
	}

	ctx := context.Background()
	tests := []test{
		{
			name:      "ok",
			input:     input,
			wantErr:   nil,
			reception: &domain.Reception{ID: uuid.New()},
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, repoerr.ErrNotFound)
				m.receptionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, nil)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:      "employee only",
			input:     dto.CreateReceptionParam{CreatorRole: domain.UserRoleModerator},
			wantErr:   domain.ErrEmployeeOnly,
			reception: nil,
			mockFn: func(test test, m receptionMocks) {

			},
		},
		{
			name:      "error previous reception is not closed",
			input:     input,
			wantErr:   domain.ErrPreviousReceptionIsNotClosed,
			reception: nil,
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, nil)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:      "error get first by status",
			input:     input,
			wantErr:   ErrReceptionOops,
			reception: nil,
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, ErrReceptionOops)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:      "error create reception",
			input:     input,
			wantErr:   ErrReceptionOops,
			reception: &domain.Reception{ID: uuid.New()},
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, repoerr.ErrNotFound)
				m.receptionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, ErrReceptionOops)
				runRepeatableReadAndUnwrap(m)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)

			useCase := NewUseCase(newDepsUseCase(mocks))

			tt.mockFn(tt, mocks)

			actual, err := useCase.Create(ctx, tt.input)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.reception, actual)
		})
	}
}

func TestReceptionUseCase_CloseLastReception(t *testing.T) {
	t.Parallel()

	type test struct {
		name            string
		input           dto.CloseLastReceptionParam
		reception       *domain.Reception
		mockFn          func(test test, m receptionMocks)
		wantErr         error
		closedReception *domain.Reception
	}

	input := dto.CloseLastReceptionParam{
		CloserRole: domain.UserRoleEmployee,
		PvzID:      uuid.New(),
	}

	ctx := context.Background()
	tests := []test{
		{
			name:            "ok",
			input:           input,
			wantErr:         nil,
			reception:       &domain.Reception{ID: uuid.New()},
			closedReception: &domain.Reception{ID: uuid.New()},
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, nil)
				m.receptionRepo.EXPECT().UpdateReceptionStatusByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.closedReception, nil)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:    "error employee only",
			input:   dto.CloseLastReceptionParam{CloserRole: domain.UserRoleModerator},
			wantErr: domain.ErrEmployeeOnly,
			mockFn: func(test test, m receptionMocks) {
			},
		},
		{
			name:    "error reception not found",
			input:   input,
			wantErr: domain.ErrReceptionNotFound,
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, repoerr.ErrNotFound)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:      "error reception already closed",
			input:     input,
			wantErr:   domain.ErrReceptionAlreadyClosed,
			reception: &domain.Reception{ID: uuid.New()},
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, nil)
				m.receptionRepo.EXPECT().UpdateReceptionStatusByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.closedReception, repoerr.ErrNotFound)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:            "error update last reception status",
			input:           input,
			wantErr:         ErrReceptionOops,
			reception:       &domain.Reception{ID: uuid.New()},
			closedReception: &domain.Reception{ID: uuid.New()},
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, nil)
				m.receptionRepo.EXPECT().UpdateReceptionStatusByID(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.closedReception, ErrReceptionOops)
				runRepeatableReadAndUnwrap(m)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)

			useCase := NewUseCase(newDepsUseCase(mocks))

			tt.mockFn(tt, mocks)

			actual, err := useCase.CloseLastReception(ctx, tt.input)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.closedReception, actual)
		})
	}
}

func TestReceptionUseCase_DeleteLastReception(t *testing.T) {
	t.Parallel()

	type test struct {
		name      string
		input     dto.DeleteLastReceptionParam
		reception *domain.Reception
		mockFn    func(test test, m receptionMocks)
		wantErr   error
	}

	input := dto.DeleteLastReceptionParam{
		DeleterRole: domain.UserRoleEmployee,
		PvzID:       uuid.New(),
	}

	ctx := context.Background()
	tests := []test{
		{
			name:      "ok",
			input:     input,
			wantErr:   nil,
			reception: &domain.Reception{ID: uuid.New()},
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, nil)
				m.productRepo.EXPECT().DeleteLastByDateTime(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:      "error get first reception by status",
			input:     input,
			wantErr:   ErrReceptionOops,
			reception: &domain.Reception{ID: uuid.New()},
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, ErrReceptionOops)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:      "error delete last product by date time",
			input:     input,
			wantErr:   ErrReceptionOops,
			reception: &domain.Reception{ID: uuid.New()},
			mockFn: func(test test, m receptionMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, nil)
				m.productRepo.EXPECT().DeleteLastByDateTime(gomock.Any(), gomock.Any()).
					Times(1).
					Return(ErrReceptionOops)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:    "error employee only",
			input:   dto.DeleteLastReceptionParam{DeleterRole: domain.UserRoleModerator},
			wantErr: domain.ErrEmployeeOnly,
			mockFn: func(test test, m receptionMocks) {
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)

			useCase := NewUseCase(newDepsUseCase(mocks))

			tt.mockFn(tt, mocks)

			err := useCase.DeleteLastProduct(ctx, tt.input)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func runRepeatableReadAndUnwrap(m receptionMocks) {
	runRepeatableRead(m)
	m.transactor.EXPECT().Unwrap(gomock.Any()).
		Times(1).
		DoAndReturn(func(err error) error {
			return err
		})
}

func runRepeatableRead(m receptionMocks) {
	m.transactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).
		Times(1).
		DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
			return transaction(ctx)
		})
}
