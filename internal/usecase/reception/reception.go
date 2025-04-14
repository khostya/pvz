//go:generate ${LOCAL_BIN}/ifacemaker -f ./reception.go -s UseCase -i Reception -p mock_reception -o ./mocks/reception.go
//go:generate ${LOCAL_BIN}/mockgen -source ./mocks/reception.go -destination=./mocks/reception.go -package=mock_reception

package reception

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
	receptionRepo receptionRepo
	productRepo   productRepo
	tm            transactionManager
}

type (
	receptionRepo interface {
		Create(ctx context.Context, reception *domain.Reception) (*domain.Reception, error)
		GetFirstByStatusAndPVZId(ctx context.Context, status domain.ReceptionStatus, pvzID uuid.UUID) (*domain.Reception, error)
		UpdateReceptionStatusByID(ctx context.Context, id uuid.UUID, status domain.ReceptionStatus) (*domain.Reception, error)
	}

	productRepo interface {
		DeleteLastByDateTimeAndReceptionID(ctx context.Context, receptionID uuid.UUID) error
	}

	DepsUseCase struct {
		ReceptionRepo      receptionRepo
		TransactionManager transactionManager
		ProductRepo        productRepo
	}

	transactionManager interface {
		RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error
		Unwrap(err error) error
	}
)

func NewUseCase(useCase DepsUseCase) *UseCase {
	return &UseCase{
		receptionRepo: useCase.ReceptionRepo,
		productRepo:   useCase.ProductRepo,
		tm:            useCase.TransactionManager,
	}
}

func (u *UseCase) Create(ctx context.Context, param dto.CreateReceptionParam) (*domain.Reception, error) {
	if param.CreatorRole != domain.UserRoleEmployee {
		return nil, domain.ErrEmployeeOnly
	}

	var reception *domain.Reception
	err := u.tm.RunRepeatableRead(ctx, func(ctx context.Context) error {
		_, err := u.receptionRepo.GetFirstByStatusAndPVZId(ctx, domain.ReceptionStatusInProgress, param.PvzID)
		if err == nil {
			return domain.ErrPreviousReceptionIsNotClosed
		}
		if !errors.Is(err, repoerr.ErrNotFound) {
			return err
		}

		reception, err = u.receptionRepo.Create(ctx, &domain.Reception{
			ID:       uuid.New(),
			DateTime: time.Now(),
			Status:   domain.ReceptionStatusInProgress,
			PvzId:    param.PvzID,
		})
		return err
	})

	return reception, u.tm.Unwrap(err)
}

func (u *UseCase) CloseLastReception(ctx context.Context, param dto.CloseLastReceptionParam) (*domain.Reception, error) {
	if param.CloserRole != domain.UserRoleEmployee {
		return nil, domain.ErrEmployeeOnly
	}

	var result *domain.Reception
	err := u.tm.RunRepeatableRead(ctx, func(ctx context.Context) error {
		reception, err := u.receptionRepo.GetFirstByStatusAndPVZId(ctx, domain.ReceptionStatusInProgress, param.PvzID)
		if errors.Is(err, repoerr.ErrNotFound) {
			return domain.ErrReceptionNotFound
		}
		if err != nil {
			return err
		}

		result, err = u.receptionRepo.UpdateReceptionStatusByID(ctx, reception.ID, domain.ReceptionStatusClose)
		if errors.Is(err, repoerr.ErrNotFound) {
			return domain.ErrReceptionAlreadyClosed
		}
		return err
	})

	return result, u.tm.Unwrap(err)
}

func (u *UseCase) DeleteLastProduct(ctx context.Context, param dto.DeleteLastReceptionParam) error {
	if param.DeleterRole != domain.UserRoleEmployee {
		return domain.ErrEmployeeOnly
	}

	err := u.tm.RunRepeatableRead(ctx, func(ctx context.Context) error {
		reception, err := u.receptionRepo.GetFirstByStatusAndPVZId(ctx, domain.ReceptionStatusInProgress, param.PvzID)
		if errors.Is(err, repoerr.ErrNotFound) {
			return domain.ErrThereIsNoInProgressReception
		}
		if err != nil {
			return err
		}

		err = u.productRepo.DeleteLastByDateTimeAndReceptionID(ctx, reception.ID)
		if errors.Is(err, repoerr.ErrNotFound) {
			return domain.ErrProductNotFound
		}
		return err
	})

	return u.tm.Unwrap(err)
}
