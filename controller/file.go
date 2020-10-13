package controller

import (
	"netdisk/util"

	"github.com/gin-gonic/gin"
)

func FileUpload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")

	if err != nil {
		util.ResponseError(c, util.CodeServerBusy)
	}
	name := fileHeader.Filename
	size := fileHeader.Size

	c.SaveUploadedFile(fileHeader, "./")
}

func FileAllInfo(c *gin.Context) {
	// 读取数据库返回

}
func FileInfo(c *gin.Context) {
	// 读取数据库 返回
}

func FileDownload(c *gin.Context) {

}
