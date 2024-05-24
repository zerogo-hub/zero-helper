package entity

import (
	"bytes"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	zerobytes "github.com/zerogo-hub/zero-helper/bytes"
	zerocodec "github.com/zerogo-hub/zero-helper/codec"
	zerologger "github.com/zerogo-hub/zero-helper/logger"
	zerotimer "github.com/zerogo-hub/zero-helper/timer"
	zeroutils "github.com/zerogo-hub/zero-helper/utils"
)

type entity struct {

	// db 数据库，如 gorm、sqlx
	readDBs []WrapReadDB
	writeDB WrapWriteDB

	// localCache 进程内缓存，本地缓存，如 bigcache、freecache
	localCache WrapCache

	// remoteCache 进程外缓存，远端缓存，如 redis
	remoteCache WrapCache

	query  QueryHandler
	update UpdateHandler
	delete DeleteHandler

	// readDBsMatchF2 读数据库数量是否符合 2 的 n 次方，求余优化
	readDBsMatchF2 bool

	// st 命中统计
	st *Stat

	// codec 编码解码，默认 msgpack
	codec zerocodec.Codec

	// twp 时间轮，定时器，用于缓存双删
	twp zerotimer.TimerWheelPool

	// timeout 给 singleflight 设置一个超时时间
	timeout time.Duration

	// notFoundExpired 数据未找到时设置短期缓存的有效期，默认 1 分钟
	notFoundExpired time.Duration

	logger zerologger.Logger
	g      singleflight.Group
	gMulti singleflight.Group
}

// New 创建一个实体管理器
func New(st *Stat, logger zerologger.Logger, codec zerocodec.Codec) Entity {
	e := &entity{
		st:              st,
		codec:           codec,
		twp:             *zerotimer.NewPool(16, 500*time.Millisecond, 120),
		timeout:         500 * time.Millisecond,
		notFoundExpired: 1 * time.Minute,
		logger:          logger,
	}

	return e
}

func (e *entity) Build() {
	e.twp.Start()

	if len(e.readDBs) == zeroutils.F2(len(e.readDBs)) {
		e.readDBsMatchF2 = true
	}
}

// Get 根据主键获取数据
func (e *entity) Get(out interface{}, id uint64) error {
	return e.get(out, id, nil)
}

// GetWithQuery 根据主键获取数据
func (e *entity) GetWithQuery(out interface{}, id uint64, query QueryHandler) error {
	return e.get(out, id, query)
}

// MGet 根据主键批量获取数据
func (e *entity) MGet(out interface{}, ids ...uint64) (*Result, error) {
	if len(ids) == 0 {
		return &Result{}, nil
	}

	key := genSingleFlightKeyMulti(ids...)

	ch := e.gMulti.DoChan(key, func() (interface{}, error) {
		result := &Result{
			IDIndexs: buildIdIndex(ids...),
			Vals:     make(map[int][]byte, len(ids)),
			Errs:     make(map[int]error, len(ids)),
		}

		missIds := make([]uint64, 0, len(ids))
		copy(missIds, ids)

		if e.localCache != nil {
			missIds = e.getMultiFromCache(e.localCache, result, ids...)
			if len(missIds) == 0 {
				// 全部命中缓存
				if e.st != nil {
					e.st.incLocalCacheHit()
				}
				return result, nil
			}

			// 存在未命中
			if e.st != nil {
				e.st.incLocalCacheMiss()
			}
		}

		if e.remoteCache != nil {
			missIds = e.getMultiFromCache(e.remoteCache, result, missIds...)
			if len(missIds) == 0 {
				// 全部命中缓存
				if e.st != nil {
					e.st.incRemoteCacheHit()
				}
				return result, nil
			}

			// 存在未命中
			if e.st != nil {
				e.st.incRemoteCacheMiss()
			}
		}

		// 以下查找出来的结果会存入缓存中

		missIdsNotCache := make([]uint64, 0, len(missIds))
		copy(missIdsNotCache, missIds)

		if len(e.readDBs) > 0 {
			// err :=
			// if err == nil {
			// 	missIds := []uint64{}

			// }
		}

		if e.query != nil {
			err := e.query(out, missIds...)
			if err == nil {
				// missIds := []uint64{}

			}
		}

		return nil, nil
	})

	select {
	case <-time.After(e.timeout):
		return nil, ErrTimeout
	case ret := <-ch:
		if ret.Err != nil {
			return nil, ret.Err
		}
		// return e.codec.Unmarshal(ret.Val.([]byte), out)
		return nil, nil
	}
}

// Update 更新
func (e *entity) Update(model interface{}, id uint64) error {
	if e.writeDB != nil {
		if err := e.writeDB.Update(model); err != nil {
			e.logger.Errorf("failed to update in db, id: %d, err: %s", id, err.Error())
			return err
		}
	}

	if e.update != nil {
		if err := e.update(model, id); err != nil {
			e.logger.Errorf("failed to update in e.update, id: %d, err: %s", id, err.Error())
		}
	}

	e.doubleDeleteCache(id)

	return nil
}

// Delete 删除
func (e *entity) Delete(model interface{}, id uint64) error {
	if e.writeDB != nil {
		if err := e.writeDB.Delete(id, model); err != nil {
			e.logger.Errorf("failed to delete in db, id: %d, err: %s", id, err.Error())
			return err
		}
	}

	if e.delete != nil {
		if err := e.delete(model, id); err != nil {
			e.logger.Errorf("failed to delete in e.delete, id: %d, err: %s", id, err.Error())
		}
	}

	e.doubleDeleteCache(id)

	return nil
}

// WithCodec 设置编码解码器
func (e *entity) WithCodec(codec zerocodec.Codec) Entity {
	e.codec = codec
	return e
}

// SetTimeout 设置超时时间
func (e *entity) WithTimeout(timeout time.Duration) Entity {
	e.timeout = timeout
	return e
}

// WithNotFoundExipred 设置短期缓存有效期
func (e *entity) WithNotFoundExipred(expired time.Duration) Entity {
	e.notFoundExpired = expired
	return e
}

func (e *entity) WithReadDB(dbs ...WrapReadDB) Entity {
	if len(e.readDBs) == 0 {
		e.readDBs = make([]WrapReadDB, 0, len(dbs))
	}
	e.readDBs = append(e.readDBs, dbs...)
	return e
}

func (e *entity) WithWriteDB(db WrapWriteDB) Entity {
	e.writeDB = db
	return e
}

func (e *entity) WithLocalCache(localCache WrapCache) Entity {
	e.localCache = localCache
	return e
}

func (e *entity) WithRemoteCache(remoteCache WrapCache) Entity {
	e.remoteCache = remoteCache
	return e
}

func (e *entity) WithCustomQueryHandler(handler QueryHandler) Entity {
	e.query = handler
	return e
}

func (e *entity) WithCustomUpdateHandler(handler UpdateHandler) Entity {
	e.update = handler
	return e
}

func (e *entity) WithCustomDeleteHandler(handler DeleteHandler) Entity {
	e.delete = handler
	return e
}

func (e *entity) get(out interface{}, id uint64, query QueryHandler) error {
	key := genSingleFlightKey(id)

	ch := e.g.DoChan(key, func() (interface{}, error) {
		var err error

		// 本地缓存
		if e.localCache != nil {
			bs, err := e.getFromLocalCache(id)
			if err == nil {
				return bs, nil
			}
			if err == ErrEmptyPlaceholder {
				return nil, ErrNotFound
			}
		}

		// 远端缓存
		if e.remoteCache != nil {
			bs, err := e.getFromRemoteCache(id)
			if err == nil {
				return bs, nil
			}
			if err == ErrEmptyPlaceholder {
				return nil, ErrNotFound
			}
		}

		// 以下查找出来的结果会存入缓存中

		if query != nil {
			// 本次单独传入的自定义查找
			err = query(out, id)

			if err != nil && e.st != nil {
				e.st.incDBFail()
			}
		} else if len(e.readDBs) > 0 {
			// 默认通过主键查找
			err = e.readDB(id).Get(id, out)

			if err != nil && e.st != nil {
				e.st.incDBFail()
			}
		} else if e.query != nil {
			// 全局自定义查找
			err = e.query(out, id)

			if err != nil && e.st != nil {
				e.st.incCustomHandlerFail()
			}
		} else {
			return nil, ErrNotFound
		}

		if err != nil {
			// 数据未找到，设置短期缓存
			e.setCacheWithNotFound(id)
			return nil, err
		}

		// 查找成功，写入缓存
		bs, err := e.codec.Marshal(out)
		if err != nil {
			e.logger.Errorf("marshal failed, id: %d, err: %s", id, err.Error())
			return nil, err
		}
		e.setCache(id, bs)

		return bs, nil
	})

	select {
	case <-time.After(e.timeout):
		return ErrTimeout
	case ret := <-ch:
		if ret.Err != nil {
			return ret.Err
		}
		return e.codec.Unmarshal(ret.Val.([]byte), out)
	}
}

// readDB 获取一条数据库读配置
//
// 数据库读缓存，求余优化
func (e *entity) readDB(id uint64) WrapReadDB {
	count := len(e.readDBs)
	if count == 1 {
		return e.readDBs[0]
	} else if count == 0 {
		return nil
	}

	if e.readDBsMatchF2 {
		return e.readDBs[id&uint64((len(e.readDBs))-1)]
	}
	return e.readDBs[id%uint64(len(e.readDBs))]
}

func (e *entity) readDBWithKey(key string) WrapReadDB {
	count := len(e.readDBs)
	if count == 1 {
		return e.readDBs[0]
	} else if count == 0 {
		return nil
	}

	v := zeroutils.ToUint64(key)

	if e.readDBsMatchF2 {
		return e.readDBs[v&uint64((len(e.readDBs))-1)]
	}
	return e.readDBs[v%uint64(len(e.readDBs))]
}

func (e *entity) getFromLocalCache(id uint64) ([]byte, error) {
	data, err := e.localCache.Get(id)
	if err != nil {
		if e.st != nil {
			e.st.incLocalCacheMiss()
		}
		return nil, err
	}

	if len(data) == 0 {
		if e.st != nil {
			e.st.incLocalCacheMiss()
		}
		return nil, e.localCache.ErrNotFound()
	}

	// 从缓存中找到数据
	if e.st != nil {
		e.st.incLocalCacheHit()
	}

	if bytes.Equal(data, emptyPlaceholder) {
		// 未命中而设置的短期缓存
		return nil, ErrEmptyPlaceholder
	}

	return data, nil
}

func (e *entity) getFromRemoteCache(id uint64) ([]byte, error) {
	data, err := e.remoteCache.Get(id)
	if err != nil {
		if e.st != nil {
			e.st.incRemoteCacheMiss()
		}
		return nil, err
	}

	if len(data) == 0 {
		if e.st != nil {
			e.st.incRemoteCacheMiss()
		}
		return nil, e.remoteCache.ErrNotFound()
	}

	// 从缓存中找到数据
	if e.st != nil {
		e.st.incRemoteCacheHit()
	}
	if bytes.Equal(data, emptyPlaceholder) {
		// 未命中而设置的短期缓存
		return nil, ErrEmptyPlaceholder
	}

	return data, nil
}

func (e *entity) setCacheWithNotFound(id uint64) {
	if e.localCache != nil {
		if err := e.localCache.Set(id, emptyPlaceholder); err != nil {
			e.logger.Errorf("set local cache with not found failed, id: %d, err: %s", id, err.Error())
			return
		}

		e.twp.AddTask(e.notFoundExpired, 1, func(t time.Time) {
			if err := e.localCache.Delete(id); err != nil {
				e.logger.Errorf("failed to delay delete in local cache, id: %d, err: %s", id, err.Error())
			}
		})
	}

	if e.remoteCache != nil {
		if err := e.remoteCache.Set(id, emptyPlaceholder); err != nil {
			e.logger.Errorf("set remote cache with not found failed, id: %d, err: %s", id, err.Error())
			return
		}

		e.twp.AddTask(e.notFoundExpired, 1, func(t time.Time) {
			if err := e.remoteCache.Delete(id); err != nil {
				e.logger.Errorf("failed to delay delete in remote cache, id: %d, err: %s", id, err.Error())
			}
		})
	}
}

func (e *entity) setCache(id uint64, bs []byte) {
	if e.localCache != nil {
		if err := e.localCache.Set(id, bs); err != nil {
			e.logger.Errorf("set local cache failed, id: %d, err: %s", id, err.Error())
		}
	}

	if e.remoteCache != nil {
		if err := e.remoteCache.Set(id, bs); err != nil {
			e.logger.Errorf("set remote cache failed, id: %d, err: %s", id, err.Error())
		}
	}
}

// doubleDeleteCache 缓存双删
func (e *entity) doubleDeleteCache(id uint64) {
	if e.localCache != nil {
		// 立即从缓存中删除
		if err := e.localCache.Delete(id); err != nil && err != e.localCache.ErrNotFound() {
			e.logger.Errorf("failed to delete in local cache, id: %d, err: %s", id, err.Error())
		}
		// 延迟从缓存中删除
		e.twp.AddTask(2*time.Second, 1, func(t time.Time) {
			if err := e.localCache.Delete(id); err != nil && err != e.localCache.ErrNotFound() {
				e.logger.Errorf("failed to delay delete in local cache, id: %d, err: %s", id, err.Error())
			}
		})
	}

	if e.remoteCache != nil {
		// 立即从缓存中删除
		if err := e.remoteCache.Delete(id); err != nil && err != e.remoteCache.ErrNotFound() {
			e.logger.Errorf("failed to delete in remote cache, id: %d, err: %s", id, err.Error())
		}
		// 延迟从缓存中删除
		e.twp.AddTask(2*time.Second, 1, func(t time.Time) {
			if err := e.remoteCache.Delete(id); err != nil && err != e.remoteCache.ErrNotFound() {
				e.logger.Errorf("failed to delay delete in remote cache, id: %d, err: %s", id, err.Error())
			}
		})
	}
}

func genSingleFlightKey(id uint64) string {
	return strconv.FormatUint(id, 10)
}

func genSingleFlightKeyMulti(ids ...uint64) string {
	if len(ids) == 1 {
		return strconv.FormatUint(ids[0], 10)
	}

	buf := buffer()
	defer releaseBuffer(buf)

	// 同一时间传入的 ids 都是相同顺序的
	for _, id := range ids {
		buf.Write(zerobytes.PutUint64(id))
	}

	return buf.String()
}

// buildIdIndex 整理 id
// [1001,1003,1005] -> {1001:0, 1003:1, 1005:2}
func buildIdIndex(ids ...uint64) map[uint64]int {
	idIndexs := make(map[uint64]int, len(ids))

	for i, id := range ids {
		idIndexs[id] = i
	}

	return idIndexs
}

func (e *entity) getMultiFromCache(cache WrapCache, result *Result, ids ...uint64) []uint64 {
	vals, _ := cache.MGet(ids...)

	missIds := []uint64{}
	for _, v := range vals {
		idx := result.IDIndexs[v.ID]
		if v.Err == nil && len(v.Val) > 0 {
			result.Vals[idx] = v.Val
		} else {
			if v.Err != nil {
				result.Errs[idx] = v.Err
			}

			missIds = append(missIds, v.ID)
		}
	}

	return missIds
}

var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{}
	bufferPool.New = func() interface{} {
		return &bytes.Buffer{}
	}
}

// buffer 从池中获取 buffer
func buffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	buff.Reset()
	return buff
}

// releaseBuffer 将 buff 放入池中
func releaseBuffer(buff *bytes.Buffer) {
	bufferPool.Put(buff)
}
