package gui

import (
	"uploader/db"
)

func InitInactiveTasksPanels(succeedTaskList []*db.UploadTaskRecord, failedTaskLists []*db.UploadTaskRecord) {
	InitSucceedTasksPanels(succeedTaskList)
	InitFailedTasksPanels(failedTaskLists)
}
