package collections

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

// Contains l 中是否包含 target
func Contains[T comparable](l []T, target T) bool {
	for _, v := range l {
		if target == v {
			return true
		}
	}
	return false
}

// Join [1,2,3] -> 1_2_3
func Join[T constraints.Ordered](l []T, sep string) string {
	l2 := make([]string, len(l))
	for i, v := range l {
		l2[i] = fmt.Sprintf("%v", v)
	}
	return strings.Join(l2, sep)
}

// Sum 计算总数 [1,2,3] -> 6
func Sum[T constraints.Integer | constraints.Float](l []T) T {
	total := T(0)
	for _, v := range l {
		total += v
	}
	return total
}

// Unique set [1,1,2,3,3] -> [1,2,3]
func Unique[T comparable](l []T) []T {
	if len(l) < 2 {
		return l
	}
	m := make(map[T]struct{}, len(l))
	for _, v := range l {
		m[v] = struct{}{}
	}
	return Keys(m)
}

// Difference 求差集
// 返回在 a 中，但是不在 b 中的集合
func Difference[T comparable](a, b []T) []T {
	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}

	diff := make([]T, 0, len(a))

	for _, v := range a {
		if _, ok := set[v]; !ok {
			diff = append(diff, v)
		}
	}

	return diff
}

func MaxCount[T comparable](a, b []T) int {
	c1 := len(a)
	c2 := len(b)

	if c1 >= c2 {
		return c1
	}
	return c2
}
