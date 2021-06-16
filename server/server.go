package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
	"uploader/db"
	"uploader/gui"
	"uploader/task"
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

	r.POST("/new", func(c *gin.Context) {
		request := &NewUploadTaskRequest{}
		err := c.BindJSON(request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ok": false,
				"msg": err,
			})
		} else {
			localPath, ok, err := gui.SelectedPath(request.IsDir, request.FileFilter)
			if err != nil {
				log.Fatal(err)
			} else {
				if ok {
					if request.TaskID == "" {
						request.TaskID = strconv.FormatInt(time.Now().Unix(), 10)
					}
					t := task.NewUploadTask(request.TaskID, localPath, request.TargetPath, request.IsDir)
					_ = db.CreateTaskRecord(task.UploadTaskToUploadTaskRecord(t))
					t.Start()

				} else {
					log.Println("file selection canceled")
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		}
	})


	r.POST("/suspend", func(c *gin.Context) {
		request := &struct {
			TaskID string `json:"taskId"`
		}{}
		err := c.BindJSON(request)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"ok": false,
				"msg": err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		}
	})

	r.POST("/cancel", func(c *gin.Context) {
		request := &struct {
			TaskID string `json:"taskId"`
		}{}
		err := c.BindJSON(request)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"ok": false,
				"msg": err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		}
	})

	r.POST("/recover", func(c *gin.Context) {
		request := &struct {
			TaskID string `json:"taskId"`
		}{}
		err := c.BindJSON(request)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"ok": false,
				"msg": err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		}
	})



	_ = r.Run(url)
}