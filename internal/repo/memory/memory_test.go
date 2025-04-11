package memory

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test(t *testing.T) {
	t.Parallel()

	cache, err := NewCache[int, string](time.Hour)
	require.NoError(t, err)

	_, ok := cache.Get(1)
	require.False(t, ok)

	cache.Put(1, "v")
	time.Sleep(time.Second)

	v, ok := cache.Get(1)
	require.True(t, ok)
	require.Equal(t, "v", v)
}

func TestTTL(t *testing.T) {
	t.Parallel()

	cache, err := NewCache[int, string](time.Nanosecond)
	require.NoError(t, err)

	cache.Put(1, "v")
	time.Sleep(time.Second)

	_, ok := cache.Get(1)
	require.False(t, ok)
}
