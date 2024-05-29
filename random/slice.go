package random

import (
	"math/rand"
)

// SliceShuffle 打乱切片，会改变传入的 l
func SliceShuffle[T any](l []T) {
	if len(l) < 2 {
		return
	}

	rand.Shuffle(len(l), func(i, j int) {
		l[i], l[j] = l[j], l[i]
	})
}
