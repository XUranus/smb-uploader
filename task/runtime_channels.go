package task

type UploadTaskRuntimeChannel struct{
	AbortChan		chan bool
	SuspendChan	 	chan bool
	ResumeChan	 	chan bool

	StatisticGUIRoutineShouldExit	chan bool
	FileCopyGUIRoutineShouldExit	chan bool
}

