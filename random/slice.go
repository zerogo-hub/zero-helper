package random

import (
	"math/rand"

	"time"
)

// SliceShuffle 打乱切片，会改变传入的 l
func SliceShuffle[T any](l []T) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(l), func(i, j int) {
		l[i], l[j] = l[j], l[i]
	})
}
