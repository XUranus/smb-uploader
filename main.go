package main

import (
	"log"
	"uploader/db"
	"uploader/gui"
	"uploader/server"
	"uploader/task"
	"uploader/util"
)




func main() {

	serverURL := util.LoadConfig()

	_ = db.ResolveUnfinishedActiveTasksStatus()

	succeedTaskList, err := db.LoadInactiveUploadTasksFromDB(db.Succeed)
	if err != nil {
		log.Fatal(err)
	}
	failedTaskList, err := db.LoadInactiveUploadTasksFromDB(db.Failed)
	if err != nil {
		log.Fatal(err)
	}

	gui.InitWindow()
	gui.InitInactiveTasksPanels(succeedTaskList, failedTaskList)
	gui.Refresh()

	server.StartServer(serverURL, true)

	gui.StartMainWindow(task.SuspendTaskIDChan, task.ResumeTaskIDChan, task.AbortTaskIDChan)
}
