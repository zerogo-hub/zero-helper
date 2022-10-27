package time_test

import (
	"testing"

	zerotime "github.com/zerogo-hub/zero-helper/time"
)

func TestMonthZero(t *testing.T) {
	at := zerotime.MonthZero(0)
	if at > 0 {
		t.Log(at)
	}
	at = zerotime.MonthZero(1)
	if at > 0 {
		t.Log(at)
	}
}
