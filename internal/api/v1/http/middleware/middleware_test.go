package middleware

import (
	"fmt"
	"github.com/google/uuid"
	mock_jwt "github.com/khostya/pvz/internal/service/jwt/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

func newManager(t *testing.T) *mock_jwt.Mockmanager {
	ctrl := gomock.NewController(t)
	return mock_jwt.NewMockmanager(ctrl)
}
func TestMiddleware_CreateValidatorMiddleware(t *testing.T) {
	mock := newManager(t)

	authenticator := NewAuthenticator(mock)
	middleware, err := CreateValidatorMiddleware(authenticator)
	require.NoError(t, err)
	require.NotNil(t, middleware)
}

func TestMiddleware_getJWTFromRequest(t *testing.T) {
	token := uuid.New().String()

	t.Run("ok", func(t *testing.T) {
		r := newRequestAuthorization(t, token)

		actual, err := getJWTFromRequest(r)
		require.NoError(t, err)
		require.Equal(t, token, actual)
	})

	t.Run("error invalid token", func(t *testing.T) {
		r := newRequestAuthorization(t, token+" "+token)

		_, err := getJWTFromRequest(r)
		require.Equal(t, ErrInvalidToken, err)
	})

	t.Run("error no auth header", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)

		_, err = getJWTFromRequest(r)
		require.Equal(t, ErrNoAuthHeader, err)
	})

	t.Run("error no Bearer prefix", func(t *testing.T) {
		r, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)
		r.Header.Set(echo.HeaderAuthorization, token)

		_, err = getJWTFromRequest(r)
		require.Equal(t, ErrInvalidToken, err)
	})
}

func newRequestAuthorization(t *testing.T, token string) *http.Request {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	r.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

	return r
}
