package utils

import (
	"hash/fnv"

	zerobytes "github.com/zerogo-hub/zero-helper/bytes"
)

// F2 返回 num 对应的 2的n次方
// 2 -> 2, 2^0
// 3 -> 4, 2^2
// 4 -> 4, 2^2
// 5 -> 8, 2^3
func F2(num int) int {
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

// ToUint64 从字符串中计算出一个数字
func ToUint64(key string) uint64 {
	h := fnv.New64a()
	h.Write(zerobytes.StringToBytes(key))
	return h.Sum64()
}
