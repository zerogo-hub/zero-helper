package time

import (
	"time"
)

// 获取某月的开始和结束时间mon为0本月,-1上月，1下月以此类推
func MonthIntervalTime(mon int) (startTime, endTime string) {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	startTime = thisMonth.AddDate(0, mon, 0).Format("2006-01-02") + " 00:00:00"
	endTime = thisMonth.AddDate(0, mon+1, -1).Format("2006-01-02") + " 23:59:59"
	return startTime, endTime
}

// MonthZero 本月第一天零点时间戳
//
// param: num 相差的月份，可以为负数, 0 表示当前月
//
// eg: MonthZero() => 1664553600
func MonthZero(num int) int64 {
	year, month, _ := time.Now().Date()

	// 当前月
	current := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	if num == 0 {
		current.Unix()
	}

	t := current.AddDate(0, num, 0)
	return t.Unix()
}
