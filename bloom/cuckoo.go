package bloom

import (
	libcuckoo "github.com/seiflotfy/cuckoofilter"

	zerobytes "github.com/zerogo-hub/zero-helper/bytes"
)

// Cuckoo 布谷鸟过滤器
type Cuckoo interface {
	// Add 添加元素
	Add(bytes []byte)

	// AddString 添加元素
	AddString(s string)

	// Del 删除元素
	Del(bytes []byte)

	// DelString 删除元素
	DelString(s string)

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

	// Count 已有数量
	Count() uint
}

type cuckoo struct {
	filter *libcuckoo.Filter
}

// NewCuckoo 容量100W大约需要内存1M
func NewCuckoo(capacity uint) Cuckoo {
	filter := libcuckoo.NewFilter(capacity)

	return &cuckoo{filter: filter}
}

// Add 添加元素
func (c *cuckoo) Add(bytes []byte) {
	c.filter.InsertUnique(bytes)
}

// AddString 添加元素
func (c *cuckoo) AddString(s string) {
	c.filter.InsertUnique(zerobytes.StringToBytes(s))
}

// Del 删除元素
func (c *cuckoo) Del(bytes []byte) {
	c.filter.Delete(bytes)
}

// DelString 删除元素
func (c *cuckoo) DelString(s string) {
	c.filter.Delete(zerobytes.StringToBytes(s))
}

// Contains 是否存在
func (c *cuckoo) Contains(bytes []byte) bool {
	return c.filter.Lookup(bytes)
}

// ContainsString 是否存在
func (c *cuckoo) ContainsString(s string) bool {
	return c.filter.Lookup(zerobytes.StringToBytes(s))
}

// ClearAll 移除所有元素
func (c *cuckoo) ClearAll() {
	c.filter.Reset()
}

// Count 已有数量
func (c *cuckoo) Count() uint {
	return c.filter.Count()
}
