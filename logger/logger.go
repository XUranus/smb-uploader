package logger

import (
	"log"
	"os"
)

var commonLogger *log.Logger = nil

func getLogger(logFileName string) *log.Logger {
	// init uploader.log
	file, err := os.OpenFile(logFileName,  os.O_CREATE | os.O_WRONLY | os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	return log.New(file, "", log.LstdFlags|log.Llongfile)
}
