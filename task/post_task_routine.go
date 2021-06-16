package task

import (
	"errors"
	"log"
	"sync"
	"time"
	"uploader/db"
	"uploader/gui"
)

/**
	PostMissionRoutine starts after all routine exit
 */
func PostMissionRoutine(lock *sync.WaitGroup, uploadTask *UploadTask) {
	lock.Wait()

	log.Println("uploadTask, enter PostMissionRoutine: ", uploadTask)

	if uploadTask.Status != UploadStatusFailed {
		uploadTask.Status = UploadStatusSucceed
	}

	uploadTask.FinishTime = time.Now().Unix()
	if uploadTask.Error != nil && uploadTask.Error.Error() == AbortError.Error() {
		uploadTask.Error = errors.New("已取消")
	}

	// sync to db
	uploadTaskRecord := UploadTaskToUploadTaskRecord(uploadTask)
	_ = db.UpdateTaskRecord(uploadTaskRecord)

	// remove for gui active list, and add for inactive list
	gui.GetMyMainWindow().Synchronize(func() {
		gui.RemoveActiveTaskPanel(uploadTask.TaskId)
		if uploadTask.Status == UploadStatusSucceed {
			_, _ = gui.AddSucceedUploadTaskPanel(&uploadTaskRecord)
		} else {
			_, _ = gui.AddFailedUploadTaskPanel(&uploadTaskRecord)
		}
	})


	uploadTask.FileCopyGUIRoutineSignal.Close()
	uploadTask.FileStatisticGUIRoutineSignal.Close()
	uploadTask.FileCopyTaskSignal.Close()
	uploadTask.FileStatisticTaskSignal.Close()
}
