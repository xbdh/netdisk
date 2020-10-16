package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"netdisk/dao/mysql"
	"netdisk/dao/redis"
	"netdisk/model"
	"netdisk/pkg/snowflake"
	"netdisk/store/oss"
	"netdisk/util"
	"os"
	"path"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type MutiPartFileInfo struct {
	UserName   string `json:"user_name"`
	FileName   string `json:"file_name"`
	FileId     uint64 `json:"file_id"`
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
	// 上传本地
	err = c.SaveUploadedFile(fileHeader, Storage+name)

	// 上传到 oss
	err = oss.UploadLocalFile("netdisk", Storage+name)
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

func FileMultiPartInit(c *gin.Context) {
	username := c.PostForm("username")
	filename := c.PostForm("filename")
	filehash := c.PostForm("filehash")
	filesize, _ := strconv.Atoi(c.PostForm("filesize"))

	sfid, _ := snowflake.GetID()

	filetotalinfo := MutiPartFileInfo{
		UserName:   username,
		FileName:   filename,
		FileId:     sfid,
		Hash:       filehash,
		Size:       int64(filesize),
		ChunkSize:  10 * 1024 * 1014,
		ChunkCount: int(math.Ceil(float64(filesize) / (10 * 1024 * 1024))),
	}

	//redis.FileMultiPartUploadInfo("mpupload", "chunkcount", filetotalinfo.ChunkCount)
	//redis.FileMultiPartUploadInfo("mpupload", "filehash", filetotalinfo.Hash)
	//redis.FileMultiPartUploadInfo("mpupload", "filesize", filetotalinfo.Size)
	//redis.FileMultiPartUploadInfo("mpupload", "fileid", filetotalinfo.FileId)

	c.JSON(http.StatusOK, filetotalinfo)

}

func FileMultiPartUpload(c *gin.Context) {
	file_id := c.Query("fileid")
	//file_hash := c.PostForm("Hash")
	index := c.Query("chunkid")
	//file_id := c.PostForm("fileid")
	////file_hash := c.PostForm("Hash")
	//index := c.PostForm("chunkid")
	//uploadfile, err := c.FormFile("file")

	fpath := "/home/chen/file/temp/" + file_id + "_" + index
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		log.Println(err)
		return
	}
	defer fd.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}
	if err != nil {
		log.Panicln(err, "打开文件错误")
		return
	}
	//c.SaveUploadedFile(uploadfile, fpath)

	//redis.FileMultiPartUploadInfo("uploadmp", "chunk"+index, 1)

	//c.ShouldBindJSON

	c.JSON(http.StatusOK, nil)

}
func FileMultiPartComplete(c *gin.Context) {
	//username := c.PostForm("username")
	filename := c.PostForm("filename")
	filehash := c.PostForm("filehash")
	fileid, _ := strconv.Atoi(c.PostForm("fileid"))
	filesize, _ := strconv.Atoi(c.PostForm("filesize"))

	fileinfo := model.FileInfo{
		FileId:   uint64(fileid),
		Size:     int64(filesize),
		Name:     filename,
		Hash:     filehash,
		Location: "",
	}

	//通知 合并 文件

	// 更新文件表，用户文件表
	if redis.FileMultiPartUploadIsComplete("upload") {
		mysql.Db.Create(&fileinfo)
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusInternalServerError, nil)
	}
	//form, _ := c.MultipartForm()

}
