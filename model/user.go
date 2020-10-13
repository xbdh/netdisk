package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId   uint64 `json:"user_id"   gorm:"type:BIGINT" `
	UserName string `json:"username"  gorm:"type:VARCHAR(64)" `
	Password string `json:"password"  gorm:"type:VARCHAR(64)"`
}
type RegisterForm struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
