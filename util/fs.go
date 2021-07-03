package util

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

// Touch creates an empty file if it doesn't exist
func Touch(fileName string) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(fileName, currentTime, currentTime)
		if err != nil {
			log.Fatal(err)
		}
	}
}
