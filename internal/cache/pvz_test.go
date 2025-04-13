package cache

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test(t *testing.T) {
	t.Parallel()

	key := uuid.New().String()
	pvz := []*domain.PVZ{{ID: uuid.New()}}
	cache, err := New[string, []*domain.PVZ](time.Hour)
	require.NoError(t, err)

	_, ok := cache.Get(key)
	require.False(t, ok)

	cache.Put(key, pvz)
	time.Sleep(time.Second)

	v, ok := cache.Get(key)
	require.True(t, ok)
	require.Equal(t, pvz, v)
}
