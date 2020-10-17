package model

import "gorm.io/gorm"

type User struct {
	UserId   uint64 `json:"user_id"   gorm:"type:BIGINT" `
	UserName string `json:"username"  gorm:"type:VARCHAR(256)" `
	Password string `json:"password"  gorm:"type:VARCHAR(256)"`
	gorm.Model
}
type RegisterForm struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
