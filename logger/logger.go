package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type MyLogger log.Logger

var CommonLogger *MyLogger
var taskLoggerMap map[string]*MyLogger

const logHome = "."

/**
	common logger is used to log common info and error
 */
func GetCommonLogger(logFileName string) (*MyLogger, error) {
	if CommonLogger != nil {
		return CommonLogger, nil
	}
	file, err := os.OpenFile(filepath.Join(logHome, logFileName),  os.O_CREATE | os.O_WRONLY | os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, err
	}
	CommonLogger = (*MyLogger)(log.New(file, "", log.LstdFlags|log.Llongfile))
	return CommonLogger, nil
}


/**
	each task logger correspond to one task, it's filename named as "${taskID}.log"
 */
func GetTaskLogger(taskId string) (*MyLogger, error) {
	if logger, ok := taskLoggerMap[taskId]; ok {
		return logger, nil
	}

	file, err := os.OpenFile(filepath.Join(logHome, taskId + ".log"),  os.O_CREATE | os.O_WRONLY | os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, err
	}
	logger := log.New(file, "", log.LstdFlags|log.Llongfile)
	taskLoggerMap[taskId] = (*MyLogger)(logger)
	return (*MyLogger)(logger), nil
}

func (logger *MyLogger) Error(function string, item interface{})  {
	(*log.Logger)(logger).Println(fmt.Sprintf("[ERROR] [%v] %v", function, item))
}

func (logger *MyLogger) Info(function string, item interface{})  {
	(*log.Logger)(logger).Println(fmt.Sprintf("[INFO] [%v] %v", function, item))
}

