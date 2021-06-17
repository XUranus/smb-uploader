package task

// global variables
var (
	SuspendTaskIDChan	chan string
	ResumeTaskIDChan	chan string
	AbortTaskIDChan		chan string
)

var taskMap map[string]*UploadTask



func init() {
	SuspendTaskIDChan = make(chan string)
	ResumeTaskIDChan = make(chan string)
	AbortTaskIDChan = make(chan string)

	taskMap = make(map[string]*UploadTask)

	go func() {
		for {
			select {
			case taskId := <-SuspendTaskIDChan:
				if uploadTask, ok := taskMap[taskId]; ok {
					uploadTask.Suspend()
				}
			case taskId := <-ResumeTaskIDChan:
				if uploadTask, ok := taskMap[taskId]; ok {
					uploadTask.Resume()
				}
			case taskId := <-AbortTaskIDChan:
				if uploadTask, ok := taskMap[taskId]; ok {
					uploadTask.Abort()
				}
			}
		}
	}()
}

