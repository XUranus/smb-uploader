package task

import (
	"errors"
	"log"
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

func (routine *RoutineSignal) Suspend() {
	routine.SuspendChan <- true
}

func (routine *RoutineSignal) Resume() {
	routine.ResumeChan <- true
}

func (routine *RoutineSignal) Abort() {
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
			log.Println("process has been suspended, drop suspend signal")
		case _ = <- routine.AbortChan:
			return true
		}
	case _ = <- routine.ResumeChan:
		log.Println("process has not been suspended, drop resume signal")
	case _ = <- routine.AbortChan:
		return true
	default:
		return false
	}
	return false
}