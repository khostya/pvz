package appctx

import (
	"context"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserID(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	_, ok := GetUserID(ctx)
	require.False(t, ok)

	ctx = SetUserID(ctx, id)

	actual, ok := GetUserID(ctx)
	require.True(t, ok)
	require.Equal(t, id, actual)
}

func TestGetIsDummy(t *testing.T) {
	ctx := context.Background()
	isDummy := true

	_, ok := GetIsDummy(ctx)
	require.False(t, ok)

	ctx = SetIsDummy(ctx, isDummy)

	actual, ok := GetIsDummy(ctx)
	require.True(t, ok)
	require.Equal(t, isDummy, actual)
}

func TestGetRole(t *testing.T) {
	ctx := context.Background()
	role := domain.UserRoleEmployee

	_, ok := GetRole(ctx)
	require.False(t, ok)

	ctx = SetRole(ctx, role)

	actual, ok := GetRole(ctx)
	require.True(t, ok)
	require.Equal(t, role, actual)
}

func TestEchoGetRole(t *testing.T) {
	e := echo.New()
	role := domain.UserRoleEmployee

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	_, ok := EchoGetRole(ctx)
	require.False(t, ok)

	EchoSetRole(ctx, role)

	actual, ok := EchoGetRole(ctx)
	require.True(t, ok)
	require.Equal(t, role, actual)
}
