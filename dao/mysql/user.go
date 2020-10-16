package mysql

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"netdisk/model"
	"netdisk/pkg/snowflake"
	"strconv"
)

const secret = "liwenzhou.com"

// encryptPassword 加密算法
func encryptPassword(data []byte) (result string) {
	h := sha1.New()
	h.Write([]byte(secret))
	fmt.Println(data)
	fmt.Println(h.Sum(data))
	fmt.Println(hex.EncodeToString(h.Sum(data)))
	return hex.EncodeToString(h.Sum(data))
}

func Register(user *model.User) (err error) {
	//sqlStr := "select count(user_id) from user where username = ?"
	//var count int64
	//err = Db.Get(&count, sqlStr, user.UserName)

	result := Db.Where("user_name = ?", user.UserName).Find(user)
	if result.Error != nil {
		return err
	}
	if result.RowsAffected > 0 {
		// 用户已存在
		return ErrorUserExit
	}
	// 生成user_id //用的snowflake不一样
	userID, err := snowflake.GetID()
	if err != nil {
		return ErrorGenIDFailed
	}
	// 生成加密密码
	password := encryptPassword([]byte(user.Password))
	// 把用户插入数据库
	//sqlStr = "insert into user(user_id, username, password) values (?,?,?)"
	err = Db.Create(&model.User{
		UserName: user.UserName,
		Password: password,
		UserId:   userID,
	}).Error

	if err != nil {
		return ErrorInsertFailed
	}
	return
}

func Login(user *model.User) (err error) {
	originPassword := user.Password // 记录一下原始密码
	//sqlStr := "select user_id, username, password from user where username = ?"
	//err = Db.Get(user, sqlStr, user.UserName)
	err = Db.Where("user_name = ?", user.UserName).Find(user).Error

	//if err != nil && err != sql.ErrNoRows {
	//	// 查询数据库出错
	//	return
	//}
	if err != nil {
		return err
	}

	password := encryptPassword([]byte(originPassword))
	if user.Password != password {
		return ErrorPasswordWrong
	}
	return
}

func GetUserByID(idStr string) (user *model.User, err error) {
	user = new(model.User)
	//sqlStr := `select user_id, username from user where user_id = ?`
	//err = Db.Get(user, sqlStr, idStr)
	id, _ := strconv.Atoi(idStr)
	Db.Where("user_id = ?", id)
	return
}
