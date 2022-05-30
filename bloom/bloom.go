package bloom

import (
	libbloom "github.com/bits-and-blooms/bloom/v3"

	zerobytes "github.com/zerogo-hub/zero-helper/bytes"
)

// Bloom 布隆过滤器接口
type Bloom interface {
	// Add 添加元素
	Add(bytes []byte)

	// AddString 添加元素
	AddString(s string)

	// Contains 是否存在
	// true: 可能存在，有误差
	// false: 一定不存在
	Contains(bytes []byte) bool

	// ContainsString 是否存在
	// true: 可能存在，有误差
	// false: 一定不存在
	ContainsString(s string) bool

	// ClearAll 移除所有元素
	ClearAll()

	// Cap 容量
	Cap() uint

	// K 要达到效果，需要使用多少次的 Hash
	K() uint
}

type bloom struct {
	filter *libbloom.BloomFilter
}

// New 创建一个布隆过滤器
// n 存放元素个数
// fp 错误率
func New(n uint, fp float64) Bloom {
	filter := libbloom.NewWithEstimates(n, fp)

	return &bloom{
		filter: filter,
	}
}

func (b *bloom) Add(bytes []byte) {
	b.filter.Add(bytes)
}

func (b *bloom) AddString(s string) {
	b.filter.Add(zerobytes.StringToBytes(s))
}

func (b *bloom) Contains(bytes []byte) bool {
	return b.filter.Test(bytes)
}

func (b *bloom) ContainsString(s string) bool {
	return b.filter.Test(zerobytes.StringToBytes(s))
}

func (b *bloom) ClearAll() {
	b.filter.ClearAll()
}

func (b *bloom) Cap() uint {
	return b.filter.Cap()
}

func (b *bloom) K() uint {
	return b.filter.K()
}
