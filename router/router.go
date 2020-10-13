package router

import (
	"netdisk/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}

	r := gin.Default()

	r.POST("/file", controller.FileUpload)
	r.GET("/files", controller.FileAllInfo)
	r.GET("/file/:file_id", controller.FileInfo)
	r.GET("/file/:file_id", controller.FileDownload)
	//r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册
	//r.POST("/signup", controller.SignUpHandler)
	// 登录
	//r.POST("/login", controller.LoginHandler)

	//r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	//	// 如果是登录的用户,判断请求头中是否有 有效的JWT  ？
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "ok",
	//	})
	//})

	//r.NoRoute(func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "404",
	//	})
	//})
	return r
}
