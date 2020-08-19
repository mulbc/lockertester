package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"math/rand"
	"path"
	"path/filepath"
	"time"

	"github.com/juju/fslock"
	log "github.com/sirupsen/logrus"
)

func main() {
	path := flag.String("path", "test.file", "path to the test file")

	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(*path))
	if err != nil {
		log.WithError(err).Fatal("Could not figure out the dir of the supplied path")
	}

	rand.Seed(time.Now().UnixNano())

	fileMutex := fslock.New(*path)

	for {
		log.WithField("path", *path).Info("Trying to lock file")
		fileMutex.Lock()
		log.WithField("path", *path).Info("Got the lock")
		content, err := ioutil.ReadFile(*path)
		if err != nil {
			log.WithError(err).Error("Error when writing to file")
		} else {
			err = checkForInconsistencies(content, dir)
			if err != nil {
				time.Sleep(time.Hour * 24)
			} else {
				log.Info("no inconsistencies yet")
			}
		}
		fileMutex.Unlock()
		n := rand.Intn(10) // n will be between 0 and 10
		log.WithField("seconds", n).Info("Sleeping")
		time.Sleep(time.Duration(n) * time.Second)
		log.Info("Wakeup")
	}
}

func checkForInconsistencies(content []byte, dir string) error {
	var compareChar byte
	for pos, char := range content {
		if compareChar == 0 {
			// First time in the for loop
			compareChar = char
			continue
		}
		if char != compareChar {
			log.
				WithField("char", string(char)).
				WithField("compareChar", string(compareChar)).
				WithField("position", pos).
				Error("Char is not CompareChar!!!")
			ioutil.WriteFile(path.Join(dir, "inconsistent.file"), content, 0777)

			return errors.New("inconsistency detected")
		}
	}
	return nil
}
