package main

import (
	"fmt"
	"log"
	"os"
	"uploader/db"
	"uploader/gui"
	"uploader/server"
	"uploader/task"
	"uploader/util"
)




func main() {

	config, err := util.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	_ = db.ResolveUnfinishedActiveTasksStatus()

	succeedTaskList, err := db.LoadInactiveUploadTasksFromDB(db.Succeed)
	if err != nil {
		log.Fatal(err)
	}
	failedTaskList, err := db.LoadInactiveUploadTasksFromDB(db.Failed)
	if err != nil {
		log.Fatal(err)
	}


	if util.PidOfPortInUse(config.ServerPort) == -1 {
		gui.InitWindow()
		gui.InitInactiveTasksPanels(succeedTaskList, failedTaskList)
		gui.Refresh()
	} else {
		gui.PopMessageBox("提示", "端口已被占用")
		os.Exit(1)
	}

	server.StartServer(fmt.Sprintf("%v:%v", config.ServerHost, config.ServerPort), true)

	gui.StartMainWindow(task.SuspendTaskIDChan, task.ResumeTaskIDChan, task.AbortTaskIDChan)
}
