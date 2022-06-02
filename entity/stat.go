package entity

import (
	"sync/atomic"
	"time"

	zerologger "github.com/zerogo-hub/zero-helper/logger"
)

// Stat 统计
type Stat struct {
	Name string

	// QueryHit 查询缓存命中次数
	QueryHit uint64
	// QueryMiss 查询缓存未命中次数
	QueryMiss uint64

	// DBFails 数据库查询失败次数
	DBFails uint64

	logger zerologger.Logger
}

// NewStat 创建一个统计对象

func NewStat(name string, logger zerologger.Logger) *Stat {
	st := &Stat{Name: name, logger: logger}
	go st.loop()
	return st
}

// IncrementQueryHit 增加查询缓存命中次数
func (st *Stat) IncrementQueryHit() {
	atomic.AddUint64(&st.QueryHit, 1)
}

// IncrementQueryMiss 增加查询未命中次数
func (st *Stat) IncrementQueryMiss() {
	atomic.AddUint64(&st.QueryMiss, 1)
}

// IncrementDbFails 增加数据库查询失败次数
func (st *Stat) IncrementDBFails() {
	atomic.AddUint64(&st.DBFails, 1)
}

// loop 进行一些统计计算
func (st *Stat) loop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// 查询
		queryHit := atomic.SwapUint64(&st.QueryHit, 0)
		queryMiss := atomic.SwapUint64(&st.QueryMiss, 0)

		dbf := atomic.SwapUint64(&st.DBFails, 0)

		st.logger.Infof("cache: %s, queryHit: %d, queryMiss: %d, db_fails: %d",
			st.Name, queryHit, queryMiss, dbf)
	}
}
