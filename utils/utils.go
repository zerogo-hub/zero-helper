package utils

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
