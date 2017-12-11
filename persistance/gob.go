package persistance

import (
	"encoding/gob"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// Encode via Gob to file
func Save(path string, object interface{}) error {
	spew.Dump(object)
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Decode Gob file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}
