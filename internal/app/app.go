package app

import (
	"context"
	"github.com/khostya/pvz/internal/api"
	"github.com/khostya/pvz/internal/config"
	"github.com/khostya/pvz/internal/lib/jwt"
	"github.com/khostya/pvz/internal/usecase"
	"github.com/khostya/pvz/pkg/hash"
	pgx "github.com/khostya/pvz/pkg/postgres"
	"github.com/khostya/pvz/pkg/postgres/transactor"
)

func StartApp(ctx context.Context, cfg config.Config) error {
	db, err := pgx.NewPool(ctx, cfg.PG.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	transactionManager := transactor.NewTransactionManager(db)
	tokenManager := jwt.NewTokenManager(jwt.ManagerDeps{
		SigningKey: cfg.Auth.SigningKey,
		AccessTTL:  cfg.Auth.AccessTokenTTL,
	})

	bcrypt := hash.NewBcryptHash(cfg.Auth.PasswordCostBcrypt)

	deps := usecase.DepsUseCase{
		Bcrypt:             bcrypt,
		TransactionManager: transactionManager,
		JwtManager:         tokenManager,
	}

	uc := usecase.NewUseCase(deps)

	return api.New(ctx, cfg, uc, tokenManager)
}
