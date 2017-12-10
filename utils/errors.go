package utils

import (
	log "github.com/sirupsen/logrus"
)

// FatalIfError throws fatal if error
func FatalIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
