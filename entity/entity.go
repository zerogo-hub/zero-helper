package entity

import (
	"bytes"
	"errors"
	"strconv"
	"time"

	"golang.org/x/sync/singleflight"

	zerocodec "github.com/zerogo-hub/zero-helper/codec"
	zerologger "github.com/zerogo-hub/zero-helper/logger"
	zerotimer "github.com/zerogo-hub/zero-helper/timer"
)

var (
	// ErrNotFound 数据未找到
	ErrNotFound = errors.New("data not found")
	// ErrEmptyPlaceholder 数据为空
	ErrEmptyPlaceholder = errors.New("empty placeholder")
	// ErrTimeout 超时
	ErrTimeout = errors.New("timeout")
)

var (
	emptyPlaceholder = []byte("*")

	g singleflight.Group
)

// QueryHandler 查询函数
type QueryHandler func(id uint64, out interface{}) error

// EntityManager 实体管理器
type EntityManager interface {
	// Get 根据主键获取数据
	Get(id uint64, out interface{}) error

	// GetWithQuery 根据主键获取数据
	GetWithQuery(id uint64, query QueryHandler, out interface{}) error

	// List 列出 key 对应的数据
	List(key string, query QueryHandler, out interface{}) error

	// Update 更新
	Update(id uint64, in interface{}) error

	// Delete 删除
	Delete(id uint64, model interface{}) error

	// SetTimeout 设置超时时间
	SetTimeout(timeout time.Duration)
}

// WrapDB 封装数据库
type WrapDB interface {
	Get(id uint64, out interface{}) error
	Update(in interface{}) error
	Delete(id uint64, model interface{}) error
	ErrNotFound() error
}

// WrapCache 封装缓存
type WrapCache interface {
	Get(key string) ([]byte, error)
	Set(key string, in []byte) error
	Delete(key string) error
	ErrNotFound() error
}

type entityManager struct {
	db     WrapDB
	cache  WrapCache
	st     *Stat
	logger zerologger.Logger
	codec  zerocodec.Codec

	// twp 时间轮，定时器
	twp zerotimer.TimerWheelPool

	// timeout 给 singleflight 设置一个超时时间
	// 如果大量相同并发请求阻塞时间过长，非常致命
	timeout time.Duration
}

// NewManager 创建一个实体管理器
func NewManager(db WrapDB, cache WrapCache, st *Stat, logger zerologger.Logger, codec zerocodec.Codec) EntityManager {
	em := &entityManager{
		db:      db,
		cache:   cache,
		st:      st,
		logger:  logger,
		codec:   codec,
		twp:     *zerotimer.NewPool(16, 500*time.Millisecond, 120),
		timeout: 500 * time.Millisecond,
	}

	em.twp.Start()

	return em
}

// Get 根据主键获取数据
func (em *entityManager) Get(id uint64, out interface{}) error {
	key := strconv.FormatUint(id, 10)
	return em.get(id, key, nil, out)
}

// GetWithQuery 根据主键获取数据
func (em *entityManager) GetWithQuery(id uint64, query QueryHandler, out interface{}) error {
	key := strconv.FormatUint(id, 10)
	return em.get(id, key, query, out)
}

func (em *entityManager) List(key string, query QueryHandler, out interface{}) error {
	return em.get(0, key, query, out)
}

// Update 更新
func (em *entityManager) Update(id uint64, in interface{}) error {
	// 先保存数据库
	if err := em.db.Update(in); err != nil {
		em.logger.Errorf("failed to update in db, id: %d, err: %s", id, err.Error())
		return err
	}

	em.doubleDeleteCache(id)

	return nil
}

// Delete 删除
func (em *entityManager) Delete(id uint64, model interface{}) error {
	// 先从数据库中删除
	if err := em.db.Delete(id, model); err != nil {
		em.logger.Errorf("failed to delete in db, id: %d, err: %s", id, err.Error())
		return err
	}

	em.doubleDeleteCache(id)

	return nil
}

// SetTimeout 设置超时时间
func (em *entityManager) SetTimeout(timeout time.Duration) {
	em.timeout = timeout
}

func (em *entityManager) get(id uint64, key string, query QueryHandler, out interface{}) error {
	ch := g.DoChan(key, func() (any, error) {
		if err := em.getFromCache(key, out); err != nil {
			if err == ErrEmptyPlaceholder {
				// 未命中缓存，且未在数据库中找到
				return nil, ErrNotFound
			} else if err != em.cache.ErrNotFound() {
				// 查找过程异常
				return nil, err
			}

			if query != nil {
				// 自定义查找
				err = query(id, out)
			} else {
				// 默认通过主键查找
				err = em.db.Get(id, out)
			}

			if err == em.db.ErrNotFound() {
				// 未从数据库中找到，设置短期缓存
				em.setCacheWithNotFound(key)
				return nil, err
			}
			if err != nil {
				// 数据库查找失败
				em.st.IncrementDBFails()
				em.logger.Errorf("query from db failed, id: %d, err: %s", id, err.Error())
				return nil, err
			}

			// 查找成功，写入缓存

			bs, err := em.codec.Marshal(out)
			if err != nil {
				em.logger.Errorf("marshal failed, id: %d, err: %s", id, err.Error())
				return nil, err
			}

			if err = em.cache.Set(key, bs); err != nil {
				em.logger.Errorf("set cache failed, id: %d, err: %s", id, err.Error())
				return nil, err
			}

			return bs, nil
		}

		return em.codec.Marshal(out)
	})

	select {
	case <-time.After(em.timeout):
		return ErrTimeout
	case ret := <-ch:
		if ret.Err != nil {
			return ret.Err
		}
		return em.codec.Unmarshal(ret.Val.([]byte), out)
	}
}

func (em *entityManager) getFromCache(key string, out interface{}) error {
	data, err := em.cache.Get(key)
	if err != nil {
		em.st.IncrementQueryMiss()
		return err
	}

	if len(data) == 0 {
		em.st.IncrementQueryMiss()
		return em.cache.ErrNotFound()
	}

	// 从缓存中找到数据
	em.st.IncrementQueryHit()
	if bytes.Compare(data, emptyPlaceholder) == 0 {
		// 未命中而设置的短期缓存
		return ErrEmptyPlaceholder
	}

	err = em.codec.Unmarshal(data, out)
	if err == nil {
		return nil
	}

	em.logger.Errorf("decode failed, key: %s, data: %v, err: %s", key, data, err.Error())

	// 删除此错误的缓存
	if err := em.cache.Delete(key); err != nil {
		em.logger.Errorf("delete invalid cache failed, id: %d, err: %s", key, data, err.Error())
	}

	return em.cache.ErrNotFound()
}

func (em *entityManager) setCacheWithNotFound(key string) {
	if err := em.cache.Set(key, emptyPlaceholder); err != nil {
		em.logger.Errorf("set cache with not found failed, key: %s, err: %s", key, err.Error())
		return
	}

	// 一分钟后删除这个短期缓存
	em.twp.AddTask(1*time.Minute, 1, func(t time.Time) {
		if err := em.cache.Delete(key); err != nil {
			em.logger.Errorf("failed to delay delete in cache, key: %s, err: %s", key, err.Error())
		}
	})
}

// doubleDeleteCache 缓存双删
func (em *entityManager) doubleDeleteCache(id uint64) {
	key := strconv.FormatUint(id, 10)

	// 立即从缓存中删除
	if err := em.cache.Delete(key); err != nil && err != em.cache.ErrNotFound() {
		em.logger.Errorf("failed to delete in cache, id: %d, err: %s", id, err.Error())
	}
	// 延迟从缓存中删除
	em.twp.AddTask(2*time.Second, 1, func(t time.Time) {
		if err := em.cache.Delete(key); err != nil && err != em.cache.ErrNotFound() {
			em.logger.Errorf("failed to delay delete in cache, id: %d, err: %s", id, err.Error())
		}
	})
}
