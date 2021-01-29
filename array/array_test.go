package array_test

import (
	"testing"

	"github.com/zerogo-hub/zero-helper/array"
)

func TestIntToString(t *testing.T) {
	if array.IntToString([]int{1, 2, 3, 4}) != "1,2,3,4" {
		t.Error("IntToString err")
	}

	if array.IntToString([]int{1, 2, 3, 4}, "+") != "1+2+3+4" {
		t.Error("IntToString err, +")
	}
}

func TestInt64ToString(t *testing.T) {
	if array.Int64ToString([]int64{1, 2, 3, 4}) != "1,2,3,4" {
		t.Error("Int64ToString err")
	}

	if array.Int64ToString([]int64{1, 2, 3, 4}, "+") != "1+2+3+4" {
		t.Error("Int64ToString err, +")
	}
}

func TestStringToString(t *testing.T) {
	if array.StringToString([]string{"1", "2", "3", "4"}) != "1,2,3,4" {
		t.Error("StringToString err")
	}

	if array.StringToString([]string{"1", "2", "3", "4"}, "+") != "1+2+3+4" {
		t.Error("StringToString err, +")
	}
}
