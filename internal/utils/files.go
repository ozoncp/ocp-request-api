package utils

import (
	"os"
	"time"
)

// WatchFile Does nothing, but constantly opens and closes file at the given path.
func WatchFile(path string) error {
	var err error
	var fp *os.File

	for err == nil {
		func() {
			fp, err = os.Open(path)
			if err != nil {
				return
			}
			defer func() {
				err = fp.Close()
			}()
		}()

		time.Sleep(time.Second)
	}
	return err
}
