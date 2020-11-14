package lib

import (
	"log"
	"os"
	"sync"
)

var (
	logger *log.Logger
	once   sync.Once
)

func GetLogger() *log.Logger {
	once.Do(func() {
		// TODO: use a log file
		// var buf bytes.Buffer
		logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)
	})
	return logger
}
