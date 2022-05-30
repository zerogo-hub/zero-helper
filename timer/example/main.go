package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	zerotimer "github.com/zerogo-hub/zero-helper/timer"
)

var (
	test = 2
)

func main() {
	if test == 1 {
		testTW()
	} else if test == 2 {
		testCron()
	}

	waitSignal()
}

func testTW() {
	tw := zerotimer.New(1*time.Second, 5)
	tw.Start()
	tw.AddTask(1*time.Second, -1, func(t time.Time) {})
}

func testCron() {
	tw := zerotimer.New(1*time.Second, 10)
	tw.Start()

	// 每天的 05:00:00 执行
	tw.AddCron(5, 0, 0, -1, func(t time.Time) {
		log.Println("cron done")
	})
	// 每天的 22:00:00 执行
	tw.AddCron(22, 0, 0, -1, func(t time.Time) {
		log.Println("cron done")
	})

	// 每周一的 05:00:00 执行
	tw.AddWeekCron(1, 5, 0, 0, -1, func(t time.Time) {
		log.Println("week cron done")
	})
	// 每周一的 22:00:00 执行
	tw.AddWeekCron(1, 22, 0, 0, -1, func(t time.Time) {
		log.Println("week cron done")
	})

	// 每月1日 05:00:00 执行
	tw.AddMonthCron(1, 5, 0, 0, -1, func(t time.Time) {
		log.Println("month cron done")
	})

	// 每年6月1日 05:00:00 执行
	tw.AddYearDayCron(6, 1, 5, 0, 0, -1, func(t time.Time) {
		log.Println("year cron done")
	})
}

// waitSignal 监听信号
func waitSignal() {
	// ctrl + c 或者 kill
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGTERM}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, sigs...)

	sig := <-ch

	signal.Stop(ch)

	log.Println(sig)
}
