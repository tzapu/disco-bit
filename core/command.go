package core

import (
	log "github.com/sirupsen/logrus"
)

// StartCommand entry point for your command
func StartCommand() {
	log.Info("Do something")

	log.Debug("Debug")
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")
	log.Fatal("Fatal")
}
