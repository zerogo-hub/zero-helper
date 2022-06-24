package entity

import (
	"strconv"
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

func (w *wrapBigcache) Get(id uint64) ([]byte, error) {
	return w.cache.Get(strconv.FormatUint(id, 10))
}

func (w *wrapBigcache) Set(id uint64, in []byte) error {
	return w.cache.Set(strconv.FormatUint(id, 10), in)
}

func (w *wrapBigcache) Delete(id uint64) error {
	return w.cache.Delete(strconv.FormatUint(id, 10))
}

func (w *wrapBigcache) ErrNotFound() error {
	return w.errNotFound
}
