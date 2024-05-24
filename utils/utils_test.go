package utils_test

import (
	"testing"

	zeroutils "github.com/zerogo-hub/zero-helper/utils"
)

func TestF2(t *testing.T) {
	r1 := zeroutils.F2(0)
	r2 := zeroutils.F2(2)
	r3 := zeroutils.F2(3)
	r4 := zeroutils.F2(4)
	r5 := zeroutils.F2(5)
	r6 := zeroutils.F2(1023)

	if r1 != 1 || r2 != 2 || r3 != 4 || r4 != 4 || r5 != 8 || r6 != 1024 {
		t.Fatal("test F2 failed")
	}
}

func TestToUint64_UppercaseLowercase(t *testing.T) {
	key := "AbCdEfGhIjKlMnOpQrStUvWxYz"
	result := zeroutils.ToUint64(key)
	if result == 0 {
		t.Errorf("ToUint64(%s) = %d; want non-zero value", key, result)
	}
}
