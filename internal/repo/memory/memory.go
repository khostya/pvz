package memory

import (
	"github.com/dgraph-io/ristretto/v2"
	"time"
)

const (
	numCounters = 1e7     // number of keys to track frequency of (10M).
	maxCost     = 1 << 30 // maximum cost of cache (1GB).
	bufferItems = 64      // number of keys per Get buffer.
)

type Cache[K ristretto.Key, V any] struct {
	cache *ristretto.Cache[K, V]
	ttl   time.Duration
}

func NewCache[K ristretto.Key, V any](ttl time.Duration) (*Cache[K, V], error) {
	cache, err := ristretto.NewCache(&ristretto.Config[K, V]{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: bufferItems,
	})
	if err != nil {
		return nil, err
	}
	return &Cache[K, V]{cache: cache, ttl: ttl}, nil
}

func (mem Cache[K, V]) Get(key K) (V, bool) {
	return mem.cache.Get(key)
}

func (mem Cache[K, V]) Put(key K, data V) {
	mem.cache.SetWithTTL(key, data, 1, mem.ttl)
}
