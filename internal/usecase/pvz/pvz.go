//go:generate ${LOCAL_BIN}/ifacemaker -f ./pvz.go -s UseCase -i Pvz -p mock_pvz -o ./mocks/pvz.go
//go:generate ${LOCAL_BIN}/mockgen -source ./mocks/pvz.go -destination=./mocks/pvz.go -package=mock_pvz

package pvz

import (
	"context"
	"github.com/khostya/pvz/internal/domain"
	"github.com/khostya/pvz/internal/dto"
)

type UseCase struct {
	pvzRepo pvzRepo
}

type (
	transactionManager interface {
		RunRepeatableRead(ctx context.Context, fx func(ctx context.Context) error) error
		RunReadCommited(ctx context.Context, fx func(ctx context.Context) error) error
		Unwrap(err error) error
	}

	pvzRepo interface {
		Create(ctx context.Context, pvz *domain.PVZ) (*domain.PVZ, error)
		GetAllPVZList(ctx context.Context) ([]*domain.PVZ, error)
		GetPVZ(ctx context.Context, param dto.GetPvzParam) ([]*domain.PVZ, error)
	}

	DepsUseCase struct {
		PvzRepo            pvzRepo
		TransactionManager transactionManager
	}
)

func NewUseCase(useCase DepsUseCase) *UseCase {
	return &UseCase{
		pvzRepo: useCase.PvzRepo,
	}
}

func (uc *UseCase) Create(ctx context.Context, param dto.CreatePvzParam) (*domain.PVZ, error) {
	if param.CreatorRole != domain.UserRoleModerator {
		return nil, domain.ErrModeratorOnly
	}
	return uc.pvzRepo.Create(ctx, &domain.PVZ{
		ID:               param.ID,
		RegistrationDate: param.RegistrationDate,
		City:             param.City,
	})
}

func (uc *UseCase) GetAllPvzList(ctx context.Context) ([]*domain.PVZ, error) {
	return uc.pvzRepo.GetAllPVZList(ctx)
}

func (uc *UseCase) GetPvz(ctx context.Context, param dto.GetPvzParam) ([]*domain.PVZ, error) {
	return uc.pvzRepo.GetPVZ(ctx, param)
}
