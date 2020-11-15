package lib

import (
	"log"
	"os"
)

var (
	logger         *log.Logger
	logDestination *os.File
)

func GetMainLogger(logfile string) *log.Logger {
	if logfile == "" {
		logDestination = os.Stdout
	} else {
		f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		logDestination = f
	}
	logger = log.New(logDestination, " ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

func CloseLogStream() {
	logDestination.Close()
}
