package entity

import (
	"time"

	bigcache "github.com/allegro/bigcache/v3"
)

// wrapBigcache 封装 bigcache
type wrapBigcache struct {
	cache       *bigcache.BigCache
	errNotFound error
}

// NewWrapCache ..
// eviction 过期时间
func NewWrapCache(eviction time.Duration) WrapCache {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(eviction))
	return &wrapBigcache{
		cache:       cache,
		errNotFound: bigcache.ErrEntryNotFound,
	}
}

func (w *wrapBigcache) Get(key string) ([]byte, error) {
	return w.cache.Get(key)
}

func (w *wrapBigcache) Set(key string, in []byte) error {
	return w.cache.Set(key, in)
}

func (w *wrapBigcache) Delete(key string) error {
	return w.cache.Delete(key)
}

func (w *wrapBigcache) ErrNotFound() error {
	return w.errNotFound
}
