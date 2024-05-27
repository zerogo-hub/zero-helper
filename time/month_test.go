package time_test

import (
	"testing"
	"time"

	zerotime "github.com/zerogo-hub/zero-helper/time"
)

func TestMonthIntervalTime(t *testing.T) {
	start, end := zerotime.MonthIntervalTime(0)
	currentYear, currentMonth, _ := time.Now().Date()
	expectedStart := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02") + " 00:00:00"
	expectedEnd := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, -1).Format("2006-01-02") + " 23:59:59"
	if start != expectedStart || end != expectedEnd {
		t.Errorf("Expected (%s, %s), but got (%s, %s)", expectedStart, expectedEnd, start, end)
	}
}

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
