//go:generate ${LOCAL_BIN}/mockgen -source ./cache.go -destination=./mocks/cache.go -package=mock_cache
package cache

type Cache[K comparable, V any] interface {
	Get(K) (V, bool)
	Put(K, V)
}
