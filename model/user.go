package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId   int64  `json:"user_id   "gorm:"type:"" `
	Name     string `json:"name      "gorm:"type:"" `
	PassWord string
}
