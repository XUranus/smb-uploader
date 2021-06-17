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

func NewTask(c *gin.Context) {
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
}


func AbortTask(c *gin.Context) {
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
		gui.AbortTaskIDChan <- request.TaskID
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	}
}


//func SuspendTask(c *gin.Context) {
//	request := &struct {
//		TaskID string `json:"taskId"`
//	}{}
//	err := c.BindJSON(request)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"ok": false,
//			"msg": err,
//		})
//	} else {
//		gui.SuspendTaskIDChan <- request.TaskID
//		c.JSON(http.StatusOK, gin.H{
//			"ok": true,
//		})
//	}
//}

//func ResumeTask(c *gin.Context) {
//	request := &struct {
//		TaskID string `json:"taskId"`
//	}{}
//	err := c.BindJSON(request)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"ok": false,
//			"msg": err,
//		})
//	} else {
//		gui.ResumeTaskIDChan <- request.TaskID
//		c.JSON(http.StatusOK, gin.H{
//			"ok": true,
//		})
//	}
//}


//func RecoverTask(c *gin.Context) {
//	request := &struct {
//		TaskID string `json:"taskId"`
//	}{}
//	err := c.BindJSON(request)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"ok": false,
//			"msg": err,
//		})
//	} else {
//		c.JSON(http.StatusOK, gin.H{
//			"ok": true,
//		})
//	}
//}