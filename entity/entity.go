package entity

import (
	"bytes"
	"errors"
	"time"

	zerocodec "github.com/zerogo-hub/zero-helper/codec"
	zerologger "github.com/zerogo-hub/zero-helper/logger"
	zerotimer "github.com/zerogo-hub/zero-helper/timer"
)

var (
	// ErrNotFound 数据未找到
	ErrNotFound = errors.New("data not found")
	// ErrEmptyPlaceholder 数据为空
	ErrEmptyPlaceholder = errors.New("empty placeholder")
)

var (
	emptyPlaceholder = []byte("*")
)

// QueryHandler 查询函数
type QueryHandler func(id uint64, out interface{}) error

// EntityManager 实体管理器
type EntityManager interface {
	// Get 根据主键获取数据
	Get(id uint64, out interface{}) error

	// GetWithQuery 根据主键获取数据
	GetWithQuery(id uint64, query QueryHandler, out interface{}) error

	// Update 更新
	Update(id uint64, in interface{}) error

	// Delete 删除
	Delete(id uint64, model interface{}) error
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
	Get(id uint64) ([]byte, error)
	Set(id uint64, in []byte) error
	Delete(id uint64) error
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
}

// NewManager 创建一个实体管理器
func NewManager(db WrapDB, cache WrapCache, st *Stat, logger zerologger.Logger, codec zerocodec.Codec) EntityManager {
	em := &entityManager{
		db:     db,
		cache:  cache,
		st:     st,
		logger: logger,
		codec:  codec,
		twp:    *zerotimer.NewPool(16, 500*time.Millisecond, 120),
	}

	em.twp.Start()

	return em
}

// Get 根据主键获取数据
func (em *entityManager) Get(id uint64, out interface{}) error {
	return em.get(id, nil, out)
}

// GetWithQuery 根据主键获取数据
func (em *entityManager) GetWithQuery(id uint64, query QueryHandler, out interface{}) error {
	return em.get(id, query, out)
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

func (em *entityManager) get(id uint64, query QueryHandler, out interface{}) error {
	var err error
	if err = em.getFromCache(id, out); err != nil {
		if err == ErrEmptyPlaceholder {
			// 未命中缓存，且未在数据库中找到
			return ErrNotFound
		} else if err != em.cache.ErrNotFound() {
			// 查找过程异常
			return err
		}

		if query != nil {
			// 自定义查找
			err = query(id, out)
		} else {
			// 默认通过主键查找
			err = em.db.Get(id, out)
		}

		if err == em.db.ErrNotFound() {
			// 未从数据库中找到
			em.setCacheWithNotFound(id)
			return err
		}
		if err != nil {
			// 数据库查找失败
			em.st.IncrementDBFails()
			em.logger.Errorf("query from db failed, id: %d, err: %s", id, err.Error())
			return err
		}

		// 查找成功，写入缓存

		bs, err := em.codec.Marshal(out)
		if err != nil {
			em.logger.Errorf("marshal failed, id: %d, err: %s", id, err.Error())
			return err
		}

		if err = em.cache.Set(id, bs); err != nil {
			em.logger.Errorf("set cache failed, id: %d, err: %s", id, err.Error())
			return err
		}

		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (em *entityManager) getFromCache(id uint64, out interface{}) error {
	data, err := em.cache.Get(id)
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

	em.logger.Errorf("decode failed, id: %d, data: %v, err: %s", id, data, err.Error())

	// 删除此错误的缓存
	if err := em.cache.Delete(id); err != nil {
		em.logger.Errorf("delete invalid cache failed, id: %d, err: %s", id, data, err.Error())
	}

	return em.cache.ErrNotFound()
}

func (em *entityManager) setCacheWithNotFound(id uint64) {
	if err := em.cache.Set(id, emptyPlaceholder); err != nil {
		em.logger.Errorf("set cache with not found failed, id: %d, err: %s", id, err.Error())
		return
	}

	// 一分钟后删除这个短期缓存
	em.twp.AddTask(1*time.Minute, 1, func(t time.Time) {
		if err := em.cache.Delete(id); err != nil {
			em.logger.Errorf("failed to delay delete in cache, id: %d, err: %s", id, err.Error())
		}
	})
}

// doubleDeleteCache 缓存双删
func (em *entityManager) doubleDeleteCache(id uint64) {
	// 立即从缓存中删除
	if err := em.cache.Delete(id); err != nil {
		em.logger.Errorf("failed to delete in cache, id: %d, err: %s", id, err.Error())
	}
	// 延迟从缓存中删除
	em.twp.AddTask(2*time.Second, 1, func(t time.Time) {
		if err := em.cache.Delete(id); err != nil {
			em.logger.Errorf("failed to delay delete in cache, id: %d, err: %s", id, err.Error())
		}
	})
}
