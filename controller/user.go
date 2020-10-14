package controller

import (
	"errors"
	"fmt"
	"net/http"
	"netdisk/dao/mysql"
	"netdisk/model"
	"netdisk/pkg/jwt"
	"netdisk/util"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignUp(c *gin.Context) {
	// 1.获取请求参数 2.校验数据有效性
	var fo model.RegisterForm
	if err := c.ShouldBindJSON(&fo); err != nil {
		util.ResponseErrorWithMsg(c, util.CodeInvalidParams, err.Error())
		return
	}
	// 3.注册用户
	err := mysql.Register(&model.User{
		UserName: fo.UserName,
		Password: fo.Password,
	})
	if errors.Is(err, mysql.ErrorUserExit) {
		util.ResponseError(c, util.CodeUserExist)
		return
	}
	if err != nil {
		zap.L().Error("mysql.Register() failed", zap.Error(err))
		util.ResponseError(c, util.CodeServerBusy)
		return
	}
	util.ResponseSuccess(c, nil)
}

func Login(c *gin.Context) {
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		zap.L().Error("invalid params", zap.Error(err))
		util.ResponseErrorWithMsg(c, util.CodeInvalidParams, err.Error())
		return
	}
	if err := mysql.Login(&u); err != nil {
		zap.L().Error("mysql.Login(&u) failed", zap.Error(err))
		util.ResponseError(c, util.CodeInvalidPassword)
		return
	}
	// 生成Token
	aToken, rToken, _ := jwt.GenToken(u.UserId)
	util.ResponseSuccess(c, gin.H{
		"accessToken":  aToken,
		"refreshToken": rToken,
		"userID":       u.UserId,
		"username":     u.UserName,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	rt := c.Query("refresh_token")
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		util.ResponseErrorWithMsg(c, util.CodeInvalidToken, "请求头缺少Auth Token")
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		util.ResponseErrorWithMsg(c, util.CodeInvalidToken, "Token格式不对")
		c.Abort()
		return
	}
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	fmt.Println(err)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}
