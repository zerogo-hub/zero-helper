// Package timer 定时器，时间轮
// 推荐使用 timer-pool.go 中的 TimerWheelPool
// 而不是直接使用 TimerWheel
package timer

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	ants "github.com/panjf2000/ants/v2"
)

var (
	// threadPool 协程池
	threadPool *ants.Pool

	// taskPool 任务池
	taskPool *sync.Pool
)

// Handler 任务执行回调函数
type Handler func(t time.Time)

// cron 类型
const (
	// cronTypeEveryDay 每日执行
	cronTypeEveryDay = 1 << iota
	// cronTypeEveryWeek 每周执行
	cronTypeEveryWeek
	// cronTypeEveryMonth 每月执行
	cronTypeEveryMonth
	// cronTypeEveryYear 每年执行
	cronTypeEveryYear
)

// Task 表示一个将要执行的任务，存储于 Slot 中
type Task struct {
	// id 此任务编号
	id uint64
	// delay 延迟时间
	delay time.Duration
	// round 在时间轮上转多少圈后才执行，当值为 0 的时候才执行
	round int
	// callback 时间到后执行的回调函数
	callback Handler
	// times 任务执行的次数，如果需要永久执行，设置为 -1
	times int

	// cron 参数
	cron                                          bool
	cronType                                      int
	month, dayOfMonth, week, hour, minute, second int
}

// nextDelay 计算下一个时刻与现在，还差多少毫秒
func (t *Task) nextDelay() time.Duration {
	if !t.cron {
		panic(fmt.Sprintf("task must be cron able, task id: %d", t.id))
	}

	switch t.cronType {
	case cronTypeEveryDay:
		return t.nextDay()
	case cronTypeEveryWeek:
		return t.nextWeek()
	case cronTypeEveryMonth:
		return t.nextMonth()
	case cronTypeEveryYear:
		return t.nextYear()
	default:
		panic(fmt.Sprintf("invalid cron type: %d", t.cronType))
	}
}

func (t *Task) nextDay() time.Duration {
	// hour, minute, second
	now := time.Now()
	year, month, dayOfMonth := now.Date()

	next := time.Date(year, month, dayOfMonth, t.hour, t.minute, t.second, 0, now.Location())

	sub := next.Sub(now)
	if sub <= 0 {
		next = next.AddDate(0, 0, 1)
		sub = next.Sub(now)
	}

	return sub
}

func (t *Task) nextWeek() time.Duration {
	// week, hour, minute, second
	now := time.Now()
	// 当前周几
	week := int(now.Weekday())
	if week == 0 {
		week = 7
	}
	year, month, dayOfMonth := now.Date()

	next := time.Date(year, month, dayOfMonth, t.hour, t.minute, t.second, 0, now.Location())

	// 匹配周几
	for t.week != week {
		next = next.AddDate(0, 1, 0)
		week = int(next.Weekday())
		if week == 0 {
			week = 7
		}
	}

	sub := next.Sub(now)
	if sub <= 0 {
		next = next.AddDate(0, 0, 7)
		sub = next.Sub(now)
	}

	return sub
}

func (t *Task) nextMonth() time.Duration {
	// dayOfMonth, hour, minute, second
	now := time.Now()
	year, month, _ := now.Date()

	next := time.Date(year, month, t.dayOfMonth, t.hour, t.minute, t.second, 0, now.Location())
	sub := next.Sub(now)
	if sub <= 0 {
		next = next.AddDate(0, 1, 0)
		sub = next.Sub(now)
	}

	return sub
}

func (t *Task) nextYear() time.Duration {
	// month, dayOfMonth, hour, minute, second
	now := time.Now()
	year, _, _ := now.Date()

	next := time.Date(year, time.Month(t.month), t.dayOfMonth, t.hour, t.minute, t.second, 0, now.Location())
	sub := next.Sub(now)
	if sub <= 0 {
		next = next.AddDate(1, 0, 0)
		sub = next.Sub(now)
	}

	return sub
}

// slot 时间槽，一个时间轮划分为多个槽，每一个槽存储着将要执行的任务
type slot struct {
	// id 槽编号
	id int
	// tasks 存储于该槽的任务
	tasks *list.List

	tw *TimerWheel
}

func (s *slot) add(task *Task) {
	s.tasks.PushBack(task)
}

func (s *slot) remove(id uint64) {
	for e := s.tasks.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.id == id {
			s.tasks.Remove(e)
			return
		}
	}
}

func (s *slot) run(t time.Time) {
	for e := s.tasks.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.round > 0 {
			task.round--
			e = e.Next()
			continue
		}

		if err := threadPool.Submit(func() {
			task.callback(t)
			taskPool.Put(task)
		}); err != nil {
			log.Println(err.Error())
		}

		next := e.Next()

		s.tasks.Remove(e)
		s.tw.afterTaskRun(task)

		e = next
	}
}

// TimerWheel 时间轮
type TimerWheel struct {
	// id 时间轮编号
	id int
	// twp 时间轮池，TimerWheelPool
	twp *TimerWheelPool

	// interval 时间轮盘精度，也就是多久跳到下一个槽中
	interval time.Duration
	ticker   *time.Ticker
	// slotNum 槽数
	slotNum int
	// slots 所有时间轮槽
	slots []*slot
	// genTaskID 用于任务 ID
	genTaskID uint64

	// pos 当前处于哪个槽中
	pos int

	// closeOnce 防止多次关闭服务
	closeOnce sync.Once
	// stop 停止执行，然后关闭
	stop bool
	// closeCh 关闭的信号
	closeCh chan bool
}

// New 创建一个时间轮
func New(interval time.Duration, slotNum int) *TimerWheel {
	slotNum = f2(slotNum)

	tw := &TimerWheel{
		interval: interval,
		slotNum:  slotNum,
		slots:    make([]*slot, 0, slotNum),
		closeCh:  make(chan bool),
	}

	for i := 0; i < slotNum; i++ {
		tw.slots = append(tw.slots, &slot{
			id:    i,
			tasks: list.New(),
			tw:    tw,
		})
	}

	return tw
}

// Start 启动时间轮
func (tw *TimerWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
}

// Close 停止并关闭时间轮
func (tw *TimerWheel) Close() {
	var once bool
	tw.closeOnce.Do(func() {
		once = true
	})

	if once {
		tw.stop = true
		tw.closeCh <- true
		close(tw.closeCh)
	}
}

// AddTask 添加一个任务
// delay 延迟时间
// times 重复执行次数，-1 表示永久重复执行
// callback 时间到之后的回调
func (tw *TimerWheel) AddTask(delay time.Duration, times int, callback Handler) *Task {
	return tw.addTask(delay, times, callback)
}

// AddCron 添加一个指定时间执行的任务
//
// 示例：每天早上 5 点执行，AddCron(5, 0, 0, -1, true, func() {})
func (tw *TimerWheel) AddCron(hour, minute, second int, times int, callback Handler) *Task {
	return tw.addCron(-1, -1, -1, hour, minute, second, cronTypeEveryDay, times, callback)
}

// AddWeekCron 添加一个每周指定时间执行的任务
// week: 1-7 表示周一到周七
//
// 示例：每周一早上 5 点执行，AddWeekCron(1, 5, 0, 0, -1, true, func() {})
func (tw *TimerWheel) AddWeekCron(week, hour, minute, second int, times int, callback Handler) *Task {
	return tw.addCron(-1, -1, week, hour, minute, second, cronTypeEveryWeek, times, callback)
}

// AddMonthCron 添加一个每月指定时间执行的任务
// dayOfMonth: 1-31
//
// 示例：每月 1 号早上 5 点执行，AddMonthCron(1, 5, 0, 0, -1, true, func() {})
func (tw *TimerWheel) AddMonthCron(dayOfMonth, hour, minute, second int, times int, callback Handler) *Task {
	return tw.addCron(-1, dayOfMonth, -1, hour, minute, second, cronTypeEveryMonth, times, callback)
}

// AddYearDayCron 添加一个具体日期执行的任务
// month: 1-12
// dayOfMonth: 1-31
//
// 示例：每年6月1日早上 5 点执行
func (tw *TimerWheel) AddYearDayCron(month, dayOfMonth, hour, minute, second int, times int, callback Handler) *Task {
	return tw.addCron(month, dayOfMonth, -1, hour, minute, second, cronTypeEveryYear, times, callback)
}

func (tw *TimerWheel) start() {
	defer func() {
		tw.Close()
	}()

	for {
		select {
		case <-tw.ticker.C:
			tw.onTick(time.Now())
		case <-tw.closeCh:
			return
		}
	}
}

func (tw *TimerWheel) onTick(t time.Time) {
	slot := tw.slots[tw.pos]
	slot.run(t)

	tw.pos++

	if tw.pos >= tw.slotNum {
		tw.pos = 0
	}
}

func (tw *TimerWheel) generateTaskID() uint64 {
	return atomic.AddUint64(&tw.genTaskID, 1)
}

func (tw *TimerWheel) addTask(delay time.Duration, times int, callback Handler) *Task {
	if delay <= 0 {
		delay = tw.interval
	}

	task := taskPool.Get().(*Task)
	task.id = tw.generateTaskID()
	task.delay = delay
	task.callback = callback
	task.times = times

	tw.put(task)

	return task
}

func (tw *TimerWheel) addCron(month, dayOfMonth, week, hour, minute, second, cronType int, times int, callback Handler) *Task {
	task := taskPool.Get().(*Task)
	task.id = tw.generateTaskID()
	task.callback = callback
	task.times = times

	task.cron = true
	task.cronType = cronType

	task.month = month
	task.dayOfMonth = dayOfMonth
	task.week = week
	task.hour = hour
	task.minute = minute
	task.second = second

	task.delay = task.nextDelay()

	tw.put(task)

	return task
}

func (tw *TimerWheel) put(task *Task) {
	round, pos := tw.calcTaskRoundAndPos(task)

	task.round = round

	slot := tw.slots[pos]
	slot.add(task)
}

func (tw *TimerWheel) afterTaskRun(task *Task) {
	if task.times == -1 || task.times > 1 {
		if task.times > 0 {
			task.times--
		}
		if task.cron {
			task.delay = task.nextDelay()
		}

		if tw.twp != nil {
			tw.twp.AddTask(task.delay, task.times, task.callback)
		} else {
			// 再次执行
			tw.addTask(task.delay, task.times, task.callback)
		}
	}
}

// calcTaskRound 计算任务需要转多少圈时间轮
func (tw *TimerWheel) calcTaskRoundAndPos(task *Task) (round, pos int) {
	delay := task.delay.Milliseconds()
	interval := tw.interval.Milliseconds()

	round = int(delay / interval / int64(tw.slotNum))
	pos = int(tw.pos+int(delay/interval)) & (tw.slotNum - 1)
	return
}

func f2(num int) int {
	if num <= 0 {
		return 1
	}

	num = num - 1
	num |= num >> 1
	num |= num >> 2
	num |= num >> 4
	num |= num >> 8
	num |= num >> 16

	return int(num + 1)
}

func init() {
	// 初始化协程池
	options := ants.Options{ExpiryDuration: 10 * time.Second, Nonblocking: true}
	threadPool, _ = ants.NewPool(256*1024, ants.WithOptions(options))

	// 初始化任务池
	taskPool = &sync.Pool{}
	taskPool.New = func() interface{} {
		return &Task{}
	}
}
