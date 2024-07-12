package main

import (
	zerologger "github.com/zerogo-hub/zero-helper/logger"
	zerotime "github.com/zerogo-hub/zero-helper/time"
)

func main() {
	caller := true
	isJSONStyle := false
	log, err := zerologger.NewLogrusLogger("testlog", "./log", caller, isJSONStyle, zerotime.Hour(7*24), zerotime.Hour(1))
	if err != nil {
		panic(err)
	}

	log.SetLevel(zerologger.INFO)
	log.Debug("Debug log")
	log.Info("Info log")
	log.Warn("Warn log")
	log.Error("Error log")
	log.Fatal("Fatal log")
}
