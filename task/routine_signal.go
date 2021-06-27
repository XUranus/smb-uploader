package task

import (
	"errors"
	"uploader/logger"
)


var AbortError = errors.New("manual abort")

type RoutineSignal struct {
	SuspendChan	chan bool
	ResumeChan	chan bool
	AbortChan	chan bool
}

func NewRoutineSignal() *RoutineSignal {
	return &RoutineSignal{
		SuspendChan: make(chan bool),
		ResumeChan: make(chan bool),
		AbortChan: make(chan bool),
	}
}

func (routine *RoutineSignal) SendSuspendSignal() {
	routine.SuspendChan <- true
}

func (routine *RoutineSignal) SendResumeSignal() {
	routine.ResumeChan <- true
}

func (routine *RoutineSignal) SendAbortSignal() {
	routine.AbortChan <- true
}

func (routine *RoutineSignal) CheckSignal() (abort bool){
	select {
	case _ = <- routine.SuspendChan:
		// block process util receive resume signal
		select {
		case _ = <- routine.ResumeChan:  // cancel block
			return false
		case _ = <- routine.SuspendChan:
			logger.CommonLogger.Info("CheckSignal", "process has been suspended, drop suspend signal")
		case _ = <- routine.AbortChan:
			return true
		}
	case _ = <- routine.ResumeChan:
		logger.CommonLogger.Info("CheckSignal", "process has not been suspended, drop resume signal")
	case _ = <- routine.AbortChan:
		return true
	default:
		return false
	}
	return false
}

func (routine *RoutineSignal) Close() {
	close(routine.SuspendChan)
	close(routine.ResumeChan)
	close(routine.AbortChan)
}