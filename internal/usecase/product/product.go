//go:generate ${LOCAL_BIN}/ifacemaker -f ./product.go -s UseCase -i Product -p mock_usecase -o ./mocks/product.go
//go:generate ${LOCAL_BIN}/mockgen -source ./mocks/product.go -destination=./mocks/product.go -package=mock_product

package product

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
	"github.com/khostya/pvz/pkg/postgres/repoerr"
	"time"
)

type UseCase struct {
	productRepo   productRepo
	receptionRepo receptionRepo
	tm            transactionManager
}

type (
	transactionManager interface {
		RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error
		RunReadCommited(ctx context.Context, fx func(ctxTX context.Context) error) error
		Unwrap(err error) error
	}

	productRepo interface {
		Create(ctx context.Context, product *domain.Product) (*domain.Product, error)
	}

	receptionRepo interface {
		GetFirstByStatus(ctx context.Context, status domain.ReceptionStatus) (*domain.Reception, error)
	}

	DepsUseCase struct {
		ProductRepo        productRepo
		TransactionManager transactionManager
		ReceptionRepo      receptionRepo
	}
)

func NewUseCase(useCase DepsUseCase) *UseCase {
	return &UseCase{
		productRepo:   useCase.ProductRepo,
		tm:            useCase.TransactionManager,
		receptionRepo: useCase.ReceptionRepo,
	}
}

func (u *UseCase) Create(ctx context.Context, param dto.CreateProductParam) (*domain.Product, error) {
	if param.CreatorRole != domain.UserRoleEmployee {
		return nil, domain.ErrEmployeeOnly
	}

	var product *domain.Product
	err := u.tm.RunRepeatableRead(ctx, func(ctx context.Context) error {
		reception, err := u.receptionRepo.GetFirstByStatus(ctx, domain.ReceptionStatusInProgress)
		if err != nil && errors.Is(err, repoerr.ErrNotFound) {
			return domain.ErrThereIsNoInProgressReception
		}

		if err != nil {
			return err
		}

		product, err = u.productRepo.Create(ctx, &domain.Product{
			ID:          uuid.New(),
			DateTime:    time.Now(),
			Type:        param.Type,
			PvzID:       param.PvzID,
			ReceptionID: reception.ID,
		})
		return err
	})
	if err != nil {
		return nil, u.tm.Unwrap(err)
	}

	return product, nil
}
