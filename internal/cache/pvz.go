package cache

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/khostya/pvz/internal/repo/memory"
	"time"
)

type (
	DefaultCache[K ristretto.Key, V any] struct {
		cache *memory.Cache[K, V]
	}
)

func New[K ristretto.Key, V any](ttl time.Duration) (*DefaultCache[K, V], error) {
	t, err := memory.NewCache[K, V](ttl)
	if err != nil {
		return nil, err
	}
	return &DefaultCache[K, V]{cache: t}, nil
}

func (h *DefaultCache[K, V]) Put(k K, v V) {
	h.cache.Put(k, v)
}

func (h *DefaultCache[K, V]) Get(k K) (V, bool) {
	return h.cache.Get(k)
}
