package utils

import (
	log "github.com/sirupsen/logrus"
)

// DebugIfError sends a debug message if error
func DebugIfError(err error) {
	if err != nil {
		log.Debug(err)
	}
}

// ErrorIfError throws fatal if error
func ErrorIfError(err error) {
	if err != nil {
		log.Error(err)
	}
}

// FatalIfError throws fatal if error
func FatalIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
