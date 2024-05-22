package utils_test

import (
	"testing"

	zeroutils "github.com/zerogo-hub/zero-helper/utils"
)

func TestF2(t *testing.T) {
	r2 := zeroutils.F2(2)
	r3 := zeroutils.F2(3)
	r4 := zeroutils.F2(4)
	r5 := zeroutils.F2(5)

	if r2 != 2 || r3 != 4 || r4 != 4 || r5 != 8 {
		t.Fatal("test F2 failed")
	}
}
