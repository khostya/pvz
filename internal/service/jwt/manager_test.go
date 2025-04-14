package jwt

import (
	"context"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	user = &domain.User{
		ID:       uuid.New(),
		Email:    "411241",
		Password: "3131",
		Role:     domain.UserRoleEmployee,
	}
)

func TestManager_GenerateToken(t *testing.T) {
	manager := NewTokenManager(ManagerDeps{
		AccessTTL:  time.Hour,
		SigningKey: "13131",
	})

	t.Run("ok dummy", func(t *testing.T) {
		token, err := manager.GenerateDummyToken(domain.UserRoleEmployee)
		require.NoError(t, err)

		assert(t, nil, domain.UserRoleEmployee, true, manager, token)
	})

	t.Run("ok token", func(t *testing.T) {
		token, err := manager.GenerateToken(user)
		require.NoError(t, err)

		assert(t, &user.ID, domain.UserRoleEmployee, false, manager, token)
	})
}

func assert(t *testing.T, userID *uuid.UUID, role domain.Role, isDummy bool, manager *Manager, token domain.Token) {
	ctx := context.Background()

	ctx, err := manager.ParseToken(ctx, string(token))
	require.NoError(t, err)

	actualUserID, ok := manager.GetUserIDFromCtx(ctx)
	if userID == nil {
		require.False(t, ok)
	} else {
		require.True(t, ok)
		require.Equal(t, *userID, actualUserID)
	}

	actualIsDummy, ok := manager.GetIsDummyFromCtx(ctx)
	require.True(t, ok)
	require.Equal(t, isDummy, actualIsDummy)

	actualRole, ok := manager.GetRoleFromCtx(ctx)
	require.True(t, ok)
	require.Equal(t, role, actualRole)
}
