package timer

import (
	"sync"
	"sync/atomic"
	"time"
)

// TimerWheelPool 时间轮池
// 推荐使用，而不是直接使用 TimerWheel
// 可以将任务分发在多个时间轮中
type TimerWheelPool struct {
	twpool []*TimerWheel
	// poolsize 池中时间轮个数，一般设置为 cpu 核心数，runtime.NumCPU()
	poolsize uint64
	incr     uint64

	// tasks Map<taskId, twid>，用于记录任务在哪个时间轮中
	tasks *sync.Map
}

// NewPool 创建一个时间轮池
func NewPool(poolsize int, interval time.Duration, slotNum int) *TimerWheelPool {
	poolsize = f2(poolsize)

	twp := &TimerWheelPool{
		twpool:   make([]*TimerWheel, 0, poolsize),
		poolsize: uint64(poolsize),
		incr:     0,
		tasks:    &sync.Map{},
	}

	for i := 0; i < poolsize; i++ {
		tw := New(interval, slotNum)
		tw.id = i
		tw.twp = twp

		twp.twpool = append(twp.twpool, tw)
	}

	return twp
}

func (twp *TimerWheelPool) Start() {
	for _, tw := range twp.twpool {
		tw.Start()
	}
}

func (twp *TimerWheelPool) Close() {
	for _, tw := range twp.twpool {
		tw.Close()
	}
}

// AddTask 添加一个任务
// delay 延迟时间
// times 重复执行次数，-1 表示永久重复执行
// callback 时间到之后的回调
func (twp *TimerWheelPool) AddTask(delay time.Duration, times int, callback Handler) *Task {
	return twp.get().AddTask(delay, times, callback)
}

// AddCron 添加一个指定时间执行的任务
//
// 示例：每天早上 5 点执行，AddCron(5, 0, 0, -1, true, func() {})
func (twp *TimerWheelPool) AddCron(hour, minute, second int, times int, callback Handler) *Task {
	return twp.get().AddCron(hour, minute, second, times, callback)
}

// AddWeekCron 添加一个每周指定时间执行的任务
// week: 1-7 表示周一到周七
//
// 示例：每周一早上 5 点执行，AddWeekCron(1, 5, 0, 0, -1, true, func() {})
func (twp *TimerWheelPool) AddWeekCron(week, hour, minute, second int, times int, callback Handler) *Task {
	return twp.get().AddWeekCron(week, hour, minute, second, times, callback)
}

// AddMonthCron 添加一个每月指定时间执行的任务
// dayOfMonth: 1-31
//
// 示例：每月 1 号早上 5 点执行，AddMonthCron(1, 5, 0, 0, -1, true, func() {})
func (twp *TimerWheelPool) AddMonthCron(dayOfMonth, hour, minute, second int, times int, callback Handler) *Task {
	return twp.get().AddMonthCron(dayOfMonth, hour, minute, second, times, callback)
}

// AddYearDayCron 添加一个具体日期执行的任务
// month: 1-12
// dayOfMonth: 1-31
//
// 示例：每年6月1日早上 5 点执行
func (twp *TimerWheelPool) AddYearDayCron(month, dayOfMonth, hour, minute, second int, times int, callback Handler) *Task {
	return twp.get().AddYearDayCron(month, dayOfMonth, hour, minute, second, times, callback)
}

// Remove 删除指定任务
func (twp *TimerWheelPool) Remove(taskID uint64) {
	twid, ok := twp.tasks.Load(taskID)
	if !ok {
		return
	}
	twp.twpool[twid.(uint64)].Remove(taskID)
}

func (twp *TimerWheelPool) get() *TimerWheel {
	idx := atomic.AddUint64(&twp.incr, 1) & (twp.poolsize - 1)
	return twp.twpool[idx]
}
