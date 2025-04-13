package api

import (
	"context"
	"github.com/khostya/pvz/internal/config"
	"github.com/khostya/pvz/internal/service/jwt"
	"github.com/khostya/pvz/internal/usecase"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err := New(ctx, config.Config{}, &usecase.UseCase{}, &jwt.Manager{})
	require.NoError(t, err)
}
