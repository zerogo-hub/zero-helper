package random_test

import (
	"testing"

	zerorandom "github.com/zerogo-hub/zero-helper/random"
)

func TestSliceShuffle(t *testing.T) {
	lint := []int{3, 2, 1, 8}
	zerorandom.SliceShuffle(lint)

	lstring := []string{"C", "D", "A", "F"}
	zerorandom.SliceShuffle(lstring)

	lnochange := []int{1}
	zerorandom.SliceShuffle(lnochange)
}
