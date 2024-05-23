package entity

import (
	"sync/atomic"
	"time"
)

// StatHandler 统计处理函数
type StatHandler func(localHit, localMiss, remoteHit, remoteMiss, dbFails, customFails uint64)

// Stat 统计
type Stat struct {
	Name string

	// localCacheHit 本地缓存命中次数
	localCacheHit uint64
	// localCacheMiss 本地缓存未命中次数
	localCacheMiss uint64

	// remoteCacheHit 查询远端缓存命中次数
	remoteCacheHit uint64
	// remoteCacheMiss 查询远端缓存未命中次数
	remoteCacheMiss uint64

	// dbFail 数据库查询失败次数
	dbFail uint64

	// customHandlerFail 自定义查询失败次数
	customHandlerFail uint64

	handler StatHandler
}

// NewStat 创建一个统计对象

func NewStat(name string, handler StatHandler) *Stat {
	st := &Stat{Name: name, handler: handler}
	if handler != nil {
		go st.loop()
	}
	return st
}

// incLocalCacheHit 增加本地缓存命中次数
func (st *Stat) incLocalCacheHit() {
	atomic.AddUint64(&st.localCacheHit, 1)
}

// incLocalCacheMiss 增加本地缓存未命中次数
func (st *Stat) incLocalCacheMiss() {
	atomic.AddUint64(&st.localCacheMiss, 1)
}

// incRemoteCacheHit 增加远端缓存命中次数
func (st *Stat) incRemoteCacheHit() {
	atomic.AddUint64(&st.remoteCacheHit, 1)
}

// incRemoteCacheMiss 增加远端缓存查询未命中次数
func (st *Stat) incRemoteCacheMiss() {
	atomic.AddUint64(&st.remoteCacheMiss, 1)
}

// incDBFail 增加数据库查询失败次数
func (st *Stat) incDBFail() {
	atomic.AddUint64(&st.dbFail, 1)
}

// incCustomHandlerFail 增加自定义查询失败次数
func (st *Stat) incCustomHandlerFail() {
	atomic.AddUint64(&st.customHandlerFail, 1)
}

// loop 进行一些统计计算
func (st *Stat) loop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		localHit := atomic.SwapUint64(&st.localCacheHit, 0)
		localMiss := atomic.SwapUint64(&st.localCacheMiss, 0)

		remoteHit := atomic.SwapUint64(&st.remoteCacheHit, 0)
		remoteMiss := atomic.SwapUint64(&st.remoteCacheMiss, 0)

		dbf := atomic.SwapUint64(&st.dbFail, 0)

		chf := atomic.SwapUint64(&st.customHandlerFail, 0)

		st.handler(localHit, localMiss, remoteHit, remoteMiss, dbf, chf)
	}
}
