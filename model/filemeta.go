package model

import "gorm.io/gorm"

type FileInfo struct {
	gorm.Model
	FileId   int64  `json:"file_id"  gorm:"type:"`
	Size     int64  `json:"size"     gorm:"type:"`
	Name     string `json:"name"     gorm:"type:"`
	Hash     string `json:"hash"     gorm:"type:"`
	Location string `json:"location" gorm:"type:"`
}
