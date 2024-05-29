package entity

import (
	"errors"

	"time"

	zerocodec "github.com/zerogo-hub/zero-helper/codec"
)

var (
	// ErrNotFound 数据未找到
	ErrNotFound = errors.New("data not found")
	// ErrEmptyPlaceholder 数据为空
	ErrEmptyPlaceholder = errors.New("empty placeholder")
	// ErrTimeout 超时
	ErrTimeout = errors.New("timeout")
	// ErrResultIndexInvalid 无效索引
	ErrResultIndexInvalid = errors.New("index out of range")
	// ErrResultIdNotFound ID 未找到
	ErrResultIdNotFound = errors.New("id not found")
	// ErrIDCantBeNull ID 不可以为空
	ErrIDCantBeNull = errors.New("id cant be null")
)

var (
	emptyPlaceholder = []byte("__z_")
)

// QueryHandler 查询函数
type QueryHandler func(out interface{}, ids ...uint64) ([]uint64, []interface{}, error)

// UpdateHandler 更新函数
type UpdateHandler func(out interface{}, id ...uint64) error

// DeleteHandler 删除函数
type DeleteHandler func(out interface{}, id ...uint64) error

// Entity 实体
type Entity interface {
	// Build 内部有一些处理
	Build()

	// Unmarshal 解码
	Unmarshal(in []byte, out interface{}) error

	// Get 根据主键获取数据
	Get(out interface{}, id uint64) error

	// GetWithQuery 根据主键获取数据
	// query 自定义查询
	GetWithQuery(out interface{}, id uint64, query QueryHandler) error

	// MGet 根据主键批量获取数据
	MGet(out interface{}, ids ...uint64) (*Result, error)

	// Update 更新数据库，更新缓存
	Update(model interface{}, id uint64) error

	// Delete 删除数据库，删除缓存
	Delete(model interface{}, id uint64) error

	// MDelete 批量删除数据库，删除缓存
	MDelete(model interface{}, ids ...uint64) error

	WithCodec(codec zerocodec.Codec) Entity
	WithTimeout(timeout time.Duration) Entity
	WithNotFoundExipred(expired time.Duration) Entity
	WithReadDB(dbs ...WrapReadDB) Entity
	WithWriteDB(db WrapWriteDB) Entity
	WithLocalCache(localCache WrapCache) Entity
	WithRemoteCache(remoteCache WrapCache) Entity
	WithCustomQueryHandler(handler QueryHandler) Entity
	WithCustomUpdateHandler(handler UpdateHandler) Entity
	WithCustomDeleteHandler(handler DeleteHandler) Entity
}

// WrapReadDB 封装读数据库
type WrapReadDB interface {
	Get(out interface{}, id uint64) error
	MGet(out interface{}, ids ...uint64) ([]uint64, []interface{}, error)
	ErrNotFound() error
}

// WrapWriteDB 封装写数据库
type WrapWriteDB interface {
	Update(in interface{}) error
	Delete(model interface{}, id uint64) error
	MDelete(model interface{}, ids ...uint64) error
	ErrNotFound() error
}

// WrapCache 封装缓存
type WrapCache interface {
	Get(id uint64) ([]byte, error)
	MGet(ids ...uint64) ([]*Value, error)
	Set(id uint64, in []byte) error
	MSet(ids []uint64, datas [][]byte) error
	Delete(id uint64) error
	MDelete(ids ...uint64) error
	ErrNotFound() error
}

type Value struct {
	ID  uint64
	Val []byte
	Err error
}
