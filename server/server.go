package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


type NewUploadTaskRequest struct {
	TaskID		string	`json:"taskId"`
	TargetPath	string	`json:"targetPath"`
	LocalPath	string	`json:"localPath"`
	IsDir		bool	`json:"isDir"`
}

func StartServer(port string) {
	r := gin.Default()

	// Serves unicode entities
	r.POST("/new", func(c *gin.Context) {
		request := &NewUploadTaskRequest{}
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

	_ = r.Run(port)
}