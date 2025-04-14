package product

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
	ErrProductOops = errors.New("product oops error")
)

type productMocks struct {
	productRepo   *mock_postgres.MockProductRepo
	transactor    *mock_transactor.MockTransactor
	receptionRepo *mock_postgres.MockReceptionRepo
}

func newProductMocks(t *testing.T) productMocks {
	ctrl := gomock.NewController(t)

	return productMocks{
		transactor:    mock_transactor.NewMockTransactor(ctrl),
		productRepo:   mock_postgres.NewMockProductRepo(ctrl),
		receptionRepo: mock_postgres.NewMockReceptionRepo(ctrl),
	}
}

func newDepsUseCase(mocks productMocks) DepsUseCase {
	return DepsUseCase{
		ProductRepo:        mocks.productRepo,
		TransactionManager: mocks.transactor,
		ReceptionRepo:      mocks.receptionRepo,
	}
}

func TestProductUseCase_Create(t *testing.T) {
	t.Parallel()

	type test struct {
		name      string
		input     dto.CreateProductParam
		product   *domain.Product
		reception *domain.Reception
		mockFn    func(ctx context.Context, test test, m productMocks)
		wantErr   error
	}

	input := dto.CreateProductParam{
		CreatorRole: domain.UserRoleEmployee,
		Type:        domain.ProductTypeShoes,
		PvzID:       uuid.New(),
	}

	ctx := context.Background()
	tests := []test{
		{
			name:      "ok",
			input:     input,
			wantErr:   nil,
			reception: &domain.Reception{ID: uuid.New()},
			product:   &domain.Product{ID: uuid.New()},
			mockFn: func(ctx context.Context, test test, m productMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, nil)
				m.productRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.product, nil)
				runRepeatableRead(m)
			},
		},
		{
			name:    "error get first reception by status not found",
			input:   input,
			wantErr: domain.ErrThereIsNoInProgressReception,
			mockFn: func(ctx context.Context, test test, m productMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, repoerr.ErrNotFound)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:    "error get first reception by status",
			input:   input,
			wantErr: ErrProductOops,
			mockFn: func(ctx context.Context, test test, m productMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, ErrProductOops)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:    "error get first reception by status",
			input:   input,
			wantErr: domain.ErrThereIsNoInProgressReception,
			mockFn: func(ctx context.Context, test test, m productMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, repoerr.ErrNotFound)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name:      "error create product",
			input:     input,
			reception: &domain.Reception{ID: uuid.New()},
			wantErr:   ErrProductOops,
			mockFn: func(ctx context.Context, test test, m productMocks) {
				m.receptionRepo.EXPECT().GetFirstByStatusAndPVZId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(test.reception, nil)
				m.productRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, ErrProductOops)
				runRepeatableReadAndUnwrap(m)
			},
		},
		{
			name: "error employee only",
			input: dto.CreateProductParam{
				CreatorRole: domain.UserRoleModerator,
				Type:        domain.ProductTypeShoes,
				PvzID:       uuid.New(),
			},
			wantErr: domain.ErrEmployeeOnly,
			mockFn: func(ctx context.Context, test test, m productMocks) {
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newProductMocks(t)

			authUseCase := NewUseCase(newDepsUseCase(mocks))

			tt.mockFn(ctx, tt, mocks)

			product, err := authUseCase.Create(ctx, tt.input)
			require.Equal(t, tt.wantErr, err)

			require.Equal(t, tt.product, product)
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
