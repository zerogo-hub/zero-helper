// Package random 雪花算法修改版
// 0 - 毫秒时间戳(41 bit) - 序列号(12 bit) - 机器 id(10 bit)
// 0 - 69年 - 4000(每秒 1000 * 4000 个 uuid) - 1000
//
// 时间回拨处理:
// 1 当相差 15 ms 之内时，等待时间追上
// 2 当超过 15 ms 时，直接替换机器 id
package random

import (
	"errors"
	"sync"

	zerotime "github.com/zerogo-hub/zero-helper/time"
)

var (
	// ErrSnowflakeWorkerID 无效的 workerID，取值范围 [0, defaultMaxWorkerID]
	ErrSnowflakeWorkerID = errors.New("bad worker id")

	// ErrSnowflakeTimeBackward 时间倒退，当前时间比上一次记录的时间还要小
	ErrSnowflakeTimeBackward = errors.New("time backward")
)

var (
	// testSnowflakeTimebackward 是否测试时间回退，测试时使用
	testSnowflakeTimebackward = false
)

const (
	// 计时起始时间，毫秒，影响 41 bit 的毫秒时间戳有效性，
	// 2^41-1 = 69 年
	// 一旦确定不可更改，默认为 2021-01-30 00:00:00.000
	defaultSnowflakeOriginTime = 1611936000000

	// 节点占用的字节数量，会影响节点数量
	defaultSnowflakeWorkerIDBits = 10

	// 毫秒内自增占用的字节数量，会影响毫秒内自增最大值
	defaultSnowflakeSequenceBits = 12
)

// SnowflakeNextWorkIDFunc ..
type SnowflakeNextWorkIDFunc func() (int, error)

// SnowflakeBackWorkIDFunc ..
type SnowflakeBackWorkIDFunc func(int) error

// Snowflake uuid 生成器
type Snowflake struct {
	// 记录上一次产生 id 的毫秒时间戳
	lastTimestamp int64

	// 当前毫秒内已生成的序列号，从 0 开始, 0 - maxSequence
	sequence uint16

	// 用来表示不同节点，这样不同节点生成的一定不同， 0 - maxWorkerID
	workerID int

	lock sync.Mutex

	// 计时起始时间，影响 41 bit 的毫秒时间戳有效性
	// 41 bit 可以存储 69 年
	originTime int64

	// 节点占用的字节数量，会影响 maxWorkerID
	workerIDBits int

	// 毫秒内自增占用的字节数量，会影响 maxSequence
	sequenceBits int

	// 节点数量
	maxWorkerID int

	// 毫秒内自增最大值
	maxSequence int

	// 当发生时间回拨时超过 15 ms，用于获取替换用的 workID，为 nil 时抛出错误
	nextWorkIDFunc SnowflakeNextWorkIDFunc

	// 当发生时间回拨时超过 15 ms，用于归还当前 workID，为 nil 时抛出错误
	backWorkIDFunc SnowflakeBackWorkIDFunc
}

// NewSnowflake 创建生成器
// workID 取值 [0,1023]
func NewSnowflake(workerID int) (*Snowflake, error) {
	return NewSnowflakeBy(workerID, defaultSnowflakeOriginTime, defaultSnowflakeWorkerIDBits, defaultSnowflakeSequenceBits, nil, nil)
}

// NewSnowflakeBy 创建生成器
func NewSnowflakeBy(workerID int, originTime int64, workerIDBits int, sequenceBits int, nextWorkIDFunc SnowflakeNextWorkIDFunc, backWorkIDFunc SnowflakeBackWorkIDFunc) (*Snowflake, error) {
	maxWorkerID := -1 ^ (-1 << workerIDBits)
	maxSequence := -1 ^ (-1 << sequenceBits)

	if workerID < 0 || workerID > maxWorkerID {
		return nil, ErrSnowflakeWorkerID
	}

	return &Snowflake{
		workerID:       workerID,
		originTime:     originTime,
		workerIDBits:   workerIDBits,
		sequenceBits:   sequenceBits,
		maxWorkerID:    maxWorkerID,
		maxSequence:    maxSequence,
		nextWorkIDFunc: nextWorkIDFunc,
		backWorkIDFunc: backWorkIDFunc,
	}, nil
}

// SetOriginTime 设置起始时间
func (snowflake *Snowflake) SetOriginTime(originTime int64) {
	snowflake.originTime = originTime
}

// SetWorkIDFunc 设置时间回拨超过 15 ms 后，用于处理 workID 的函数
func (snowflake *Snowflake) SetWorkIDFunc(nextWorkIDFunc SnowflakeNextWorkIDFunc, backWorkIDFunc SnowflakeBackWorkIDFunc) {
	snowflake.nextWorkIDFunc = nextWorkIDFunc
	snowflake.backWorkIDFunc = backWorkIDFunc
}

func (snowflake *Snowflake) setWorkID(workerID int) {
	snowflake.workerID = workerID
}

// UUID 获取 uuid，线程安全
func (snowflake *Snowflake) UUID() (uint64, error) {
	snowflake.lock.Lock()
	defer snowflake.lock.Unlock()

	return snowflake.generateUUID()
}

// UnsafeUUID 获取 uuid，非线程安全
func (snowflake *Snowflake) UnsafeUUID() (uint64, error) {
	return snowflake.generateUUID()
}

func (snowflake *Snowflake) generateUUID() (uint64, error) {
	t := snowflakeNow() - snowflake.originTime
	// 时间回拨，可能会使得产生的 uuid 重复
	if t < snowflake.lastTimestamp {
		if snowflake.lastTimestamp-t <= 15 {
			// 短期内(15 ms) 等待服务器时间追上
			t = snowflake.wait()
		} else if snowflake.nextWorkIDFunc != nil && snowflake.backWorkIDFunc != nil {
			// 使用替换未使用过的 workID 的方式来生成唯一的 uuid

			// 获取一个新的 workID
			nextWorkID, err := snowflake.nextWorkIDFunc()
			if err != nil {
				return 0, ErrSnowflakeTimeBackward
			}

			curWorkID := snowflake.workerID

			// 归还当前的 workID，注意避免该 workID 立即被另一个节点获取，一般插到队列的尾部
			if err := snowflake.backWorkIDFunc(curWorkID); err != nil {
				return 0, ErrSnowflakeTimeBackward
			}

			snowflake.setWorkID(nextWorkID)
		} else {
			return 0, ErrSnowflakeTimeBackward
		}
	}

	if t == snowflake.lastTimestamp {
		// 上一次与当前都在同一个毫秒内，递增数量
		snowflake.sequence = (snowflake.sequence + 1) & uint16(snowflake.maxSequence)

		// 如果已经超出当前毫秒可以记录的范围 maxSequence
		// 1000 & 0111 => 0
		if snowflake.sequence == 0 {
			t = snowflake.wait()
		}
	} else {
		snowflake.sequence = 0
	}

	snowflake.lastTimestamp = t

	// 将相关数据封装成 uint64
	v1 := snowflake.workerIDBits + snowflake.sequenceBits
	v2 := snowflake.sequenceBits

	uuid := ((uint64(snowflake.lastTimestamp) << v1) | (uint64(snowflake.workerID << v2)) | (uint64(snowflake.sequence)))
	return uuid, nil
}

// wait 等到当前时间 > 上一次的时间
func (snowflake *Snowflake) wait() int64 {
	t := snowflakeNow()

	for t <= snowflake.lastTimestamp {
		zerotime.SleepMircosecond(100)
		t = snowflakeNow()
	}

	return t
}

func snowflakeNow() int64 {
	if !testSnowflakeTimebackward {
		return zerotime.MS()
	}
	return zerotime.MS() - 10000000
}

// SetSnowflakeTestTimebackward 测试时间回退
func SetSnowflakeTestTimebackward(able bool) {
	testSnowflakeTimebackward = able
}
