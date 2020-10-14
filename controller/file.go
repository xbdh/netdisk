package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"netdisk/dao/mysql"
	"netdisk/model"
	"netdisk/pkg/snowflake"
	"netdisk/util"
	"os"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type MutiPartFileInfo struct {
	UserName   string `json:"user_name"`
	FileName   string `json:"file_name"`
	Hash       string `json:"hash"`
	Size       int64  `json:"size"`
	ChunkSize  int    `json:"chunk_size"`
	ChunkCount int    `json:"chunk_count"`
}

const Storage = "/home/chen/file/netdisk/"

func FileUpload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	//hash := c.Query("hash")
	if err != nil {
		util.ResponseError(c, util.CodeServerBusy)
	}
	name := fileHeader.Filename
	size := fileHeader.Size
	uuid, _ := snowflake.GetID()
	f, _ := fileHeader.Open()
	hash, _ := util.GetFileHash(f)

	fileinfo := model.FileInfo{
		FileId: uuid,
		Size:   size,
		Name:   name,
		Hash:   hash,
	}

	err = mysql.Db.Create(&fileinfo).Error

	if err != nil {
		fmt.Printf("数据上传数据库失败 55555%v\n", err)
		return
	}

	err = c.SaveUploadedFile(fileHeader, Storage+name)
	if err != nil {
		fmt.Println("保存失败", err)
		return
	}
	c.JSON(http.StatusOK, fileinfo)
}

func FileAllInfo(c *gin.Context) {
	fileinfos := []model.FileInfo{}
	// 读取数据库返回

	err := mysql.Db.Find(&fileinfos).Error
	if err != nil {
		zap.L().Error("查找失败", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, fileinfos)

}
func FileInfo(c *gin.Context) {
	file_id := c.Param("file_id")

	fileinfo := model.FileInfo{}
	err := mysql.Db.Where("file_id = ?", file_id).Find(&fileinfo).Error

	if err != nil {
		zap.L().Error("查找信息失败", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, fileinfo)
}

func FileDownload(c *gin.Context) {
	file_id := c.Param("file_id")

	fileinfo := model.FileInfo{}
	err := mysql.Db.Where("file_id = ?", file_id).Find(&fileinfo).Error

	if err != nil {
		zap.L().Error("查找信息失败", zap.Error(err))
		return
	}
	name := fileinfo.Name
	f, err := os.Open(Storage + name)
	if err != nil {
		zap.L().Error("打开文件错误", zap.Error(err))
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)

		zap.L().Error("读取文件错误", zap.Error(err))
		return
	}
	c.Writer.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	c.Writer.Header().Set("content-disposition", "attachment;filename="+name)
	c.Writer.Write(data)
}

func FileMultiPartUpload(c *gin.Context) {

	//c.ShouldBindJSON

}
func FileMultiPartInit(c *gin.Context) {

}
func FileMultiPartComplete(c *gin.Context) {

}
