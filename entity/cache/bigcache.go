package cache

import (
	"context"
	"errors"
	"strconv"
	"time"

	bigcache "github.com/allegro/bigcache/v3"
	zeroentity "github.com/zerogo-hub/zero-helper/entity"
)

// wrapBigcache 封装 bigcache
type wrapBigcache struct {
	cache       *bigcache.BigCache
	errNotFound error
}

// NewBigCache ..
//
// eviction 过期时间
//
// 缺点: 无法单独对指定的 key 单独设置过期时间
func NewBigCache(eviction time.Duration) zeroentity.WrapCache {
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(eviction))
	return &wrapBigcache{
		cache:       cache,
		errNotFound: bigcache.ErrEntryNotFound,
	}
}

func (w *wrapBigcache) Get(id uint64) ([]byte, error) {
	return w.cache.Get(strconv.FormatUint(id, 10))
}

func (w *wrapBigcache) MGet(ids ...uint64) ([]*zeroentity.Value, error) {
	// 不支持批量获取，改为遍历获取
	results := make([]*zeroentity.Value, 0, len(ids))
	for _, id := range ids {
		val, err := w.Get(id)
		results = append(results, &zeroentity.Value{ID: id, Val: val, Err: err})
	}

	return results, nil
}

func (w *wrapBigcache) Set(id uint64, in []byte) error {
	return w.cache.Set(strconv.FormatUint(id, 10), in)
}

func (w *wrapBigcache) MSet(ids []uint64, datas [][]byte) error {
	if len(ids) != len(datas) {
		return errors.New("invalid length")
	}

	for idx, id := range ids {
		if err := w.Set(id, datas[idx]); err != nil {
			return err
		}
	}

	return nil
}

func (w *wrapBigcache) Delete(id uint64) error {
	return w.cache.Delete(strconv.FormatUint(id, 10))
}

func (w *wrapBigcache) MDelete(ids ...uint64) error {
	for _, id := range ids {
		if err := w.Delete(id); err != nil && err != w.errNotFound {
			return err
		}
	}
	return nil
}

func (w *wrapBigcache) ErrNotFound() error {
	return w.errNotFound
}
