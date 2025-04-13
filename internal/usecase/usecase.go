package usecase

import (
	"github.com/khostya/pvz/internal/repo/postgres"
	"github.com/khostya/pvz/internal/service/jwt"
	"github.com/khostya/pvz/internal/usecase/auth"
	"github.com/khostya/pvz/internal/usecase/product"
	"github.com/khostya/pvz/internal/usecase/pvz"
	"github.com/khostya/pvz/internal/usecase/reception"
	"github.com/khostya/pvz/pkg/hash"
	"github.com/khostya/pvz/pkg/postgres/transactor"
)

type (
	DepsUseCase struct {
		TransactionManager *transactor.TransactionManager
		JwtManager         *jwt.Manager
		Bcrypt             *hash.BcryptHash
	}

	UseCase struct {
		Auth      *auth.UseCase
		Pvz       *pvz.UseCase
		Reception *reception.UseCase
		Product   *product.UseCase
	}
)

func NewUseCase(deps DepsUseCase) *UseCase {
	repo := postgres.NewRepositories(deps.TransactionManager)

	return &UseCase{
		Auth: auth.NewAuthUseCase(auth.AuthDepsUseCase{
			UserRepo:       repo.UserRepo,
			JwtManager:     deps.JwtManager,
			PasswordHasher: deps.Bcrypt,
		}),
		Pvz: pvz.NewUseCase(pvz.DepsUseCase{
			PvzRepo:            repo.PvzRepo,
			TransactionManager: deps.TransactionManager,
		}),
		Reception: reception.NewUseCase(reception.DepsUseCase{
			TransactionManager: deps.TransactionManager,
			ReceptionRepo:      repo.ReceptionRepo,
			ProductRepo:        repo.ProductRepo,
		}),
		Product: product.NewUseCase(product.DepsUseCase{
			TransactionManager: deps.TransactionManager,
			ProductRepo:        repo.ProductRepo,
			ReceptionRepo:      repo.ReceptionRepo,
		}),
	}
}
