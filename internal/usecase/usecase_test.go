package usecase

import (
	"github.com/khostya/pvz/internal/lib/jwt"
	"github.com/khostya/pvz/pkg/hash"
	"github.com/khostya/pvz/pkg/postgres/transactor"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewUseCase(t *testing.T) {
	deps := DepsUseCase{
		TransactionManager: &transactor.TransactionManager{},
		JwtManager:         &jwt.Manager{},
		Bcrypt:             hash.NewBcryptHash(3),
	}

	uc := NewUseCase(deps)
	require.NotNil(t, uc)
	require.NotNil(t, uc.Pvz)
	require.NotNil(t, uc.Auth)
	require.NotNil(t, uc.Product)
	require.NotNil(t, uc.Reception)
}
