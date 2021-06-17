package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"uploader/db"
	"uploader/gui"
	"uploader/server"
	"uploader/task"
	"uploader/util"
)




func main() {
	// abs home path of /SmbUploader
	homePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// load config from .\config.ini
	config, err := util.LoadConfig(homePath)
	if err != nil {
		log.Fatal(err)
	}

	// invoked by protocol url
	if len(os.Args) > 2 {
		_ = util.ExtractAndParseCommandLineArg1(config.Protocol, os.Args[1])
	}


	// resolve and load persisted task from .\data.db
	db.InitDbPath(homePath)
	_ = db.ResolveUnfinishedActiveTasksStatus()

	succeedTaskList, err := db.LoadInactiveUploadTasksFromDB(db.Succeed)
	if err != nil {
		log.Fatal(err)
	}
	failedTaskList, err := db.LoadInactiveUploadTasksFromDB(db.Failed)
	if err != nil {
		log.Fatal(err)
	}


	// start GUI
	if util.PidOfPortInUse(config.ServerPort) == -1 {
		gui.InitResourcePath(homePath)
		gui.InitWindow()
		gui.InitInactiveTasksPanels(succeedTaskList, failedTaskList)
		gui.RefreshMainWindow()
	} else {
		gui.PopMessageBox("提示", "端口已被占用")
		os.Exit(1)
	}

	// start HTTP server
	server.StartServer(fmt.Sprintf("%v:%v", config.ServerHost, config.ServerPort), true)

	// GUI must run in main thread
	gui.StartMainWindowBlock(task.SuspendTaskIDChan, task.ResumeTaskIDChan, task.AbortTaskIDChan)
}
