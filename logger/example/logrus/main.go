package main

import (
	"github.com/zerogo-hub/zero-helper/logger"
	"github.com/zerogo-hub/zero-helper/time"
)

func main() {
	log, err := logger.NewLogrusLogger("testlog", "./log", true, true, time.Hour(7*24), time.Hour(1))
	if err != nil {
		panic(err)
	}

	log.Debug("Debug log")
	log.Info("Info log")
	log.Warn("Warn log")
	log.Error("Error log")
	log.Fatal("Fatal log")
}
