package task

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"uploader/db"
	"uploader/gui"
	"uploader/logger"
)

/**
	PostMissionRoutine starts after all routine exit
 */
func PostMissionRoutine(lock *sync.WaitGroup, uploadTask *UploadTask) {
	lock.Wait()

	logger.CommonLogger.Info("PostMissionRoutine", fmt.Sprintf("uploadTask, enter PostMissionRoutine: %v", uploadTask))

	if uploadTask.Status != UploadStatusFailed {
		uploadTask.Status = UploadStatusSucceed
	}

	uploadTask.FinishTime = time.Now().Unix()
	if uploadTask.Error != nil && uploadTask.Error.Error() == AbortError.Error() {
		uploadTask.Error = errors.New("已取消")
	}

	// sync to db
	uploadTaskRecord := CovertUploadTaskToUploadTaskRecord(uploadTask)
	_ = db.UpdateTaskRecord(uploadTaskRecord)

	// remove for gui active list, and add for inactive list
	gui.GetMyMainWindow().Synchronize(func() {

		gui.RemoveActiveTaskPanel(uploadTask.TaskId)

		if uploadTask.Status == UploadStatusSucceed {
			_, _ = gui.AddSucceedUploadTaskPanel(&uploadTaskRecord)
			gui.ShowCustomNotify("上传完成", uploadTask.LocalPath)

		} else {
			_, _ = gui.AddFailedUploadTaskPanel(&uploadTaskRecord)
			gui.ShowCustomNotify("上传失败",uploadTask.LocalPath )
		}
	})


	uploadTask.FileCopyGUIRoutineSignal.Close()
	uploadTask.FileStatisticGUIRoutineSignal.Close()
	uploadTask.FileCopyTaskSignal.Close()
	uploadTask.FileStatisticTaskSignal.Close()
}
