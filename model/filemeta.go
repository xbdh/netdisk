package model

import "gorm.io/gorm"

type FileInfo struct {
	gorm.Model
	FileId   uint64 `json:"file_id"  gorm:"type:BIGINT"`
	Size     int64  `json:"size"     gorm:"type:BIGINT"`
	Name     string `json:"name"     gorm:"type:VARCHAR(64)"`
	Hash     string `json:"hash"     gorm:"type:VARCHAR(64)"`
	Location string `json:"location" gorm:"type:VARCHAR(64)"`
}
