package main

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/juju/fslock"
	log "github.com/sirupsen/logrus"
)

func main() {
	charToWrite := flag.String("char", "a", "a character that is used to fill the file")
	fileSize := flag.Int("size", 1000, "amount of chars to write to the file")
	path := flag.String("path", "test.file", "path to the test file")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	fileMutex := fslock.New(*path)

	content := make([]byte, len(*charToWrite)**fileSize)
	bp := 0
	for i := 0; i < *fileSize; i++ {
		bp += copy(content[bp:], *charToWrite)
	}

	for {
		log.WithField("path", *path).Info("Trying to lock file")
		fileMutex.Lock()
		log.WithField("path", *path).Info("Got the lock")
		err := ioutil.WriteFile(*path, content, 0777)
		fileMutex.Unlock()
		if err != nil {
			log.WithError(err).Error("Error when writing to file")
			continue
		}
		n := rand.Intn(10) // n will be between 0 and 10
		log.WithField("seconds", n).Info("Sleeping")
		time.Sleep(time.Duration(n) * time.Second)
		log.Info("Wakeup")
	}
}
