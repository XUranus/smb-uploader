package server

import (
	"github.com/gin-gonic/gin"
)


type NewUploadTaskRequest struct {
	TaskID		string	`json:"taskId"`
	TargetPath	string	`json:"targetPath"`
	IsDir		bool	`json:"isDir"`
	FileFilter	string	`json:"fileFilter"`
}


func StartServer(url string, async bool) {
	if async {
		go StartServerBlock(url)
	} else {
		StartServerBlock(url)
	}
}

func StartServerBlock(url string) {
	r := gin.Default()

	// create a new upload task from HTTP Request
	r.POST("/new", NewTask)

	// judge whether program is running
	r.GET("/live", IsAlive)

	// abort a running task
	r.POST("/abort", AbortTask)

	//TODO:: suspend a running task
	//r.POST("/suspend", SuspendTask)

	//TODO:: resume a running task
	//r.POST("/resume", ResumeTask)

	//TODO:: recover a running task
	//r.POST("/recover", RecoverTask)

	_ = r.Run(url)
}