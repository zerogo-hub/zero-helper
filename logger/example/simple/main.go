package main

import (
	"github.com/zerogo-hub/zero-helper/logger"
)

func main() {
	log := logger.NewSampleLogger()

	log.Debug("Debug log")
	log.Info("Info log")
	log.Warn("Warn log")
	log.Error("Error log")
	log.Fatal("Fatal log")
}
