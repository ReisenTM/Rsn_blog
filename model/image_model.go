package model

import "fmt"

// ImageModel 图片表
type ImageModel struct {
	Model
	Filename string `gorm:"size:32" json:"filename"`
	Path     string `gorm:"size:256" json:"path"`
	Size     int64  `json:"size"`
	Hash     string `gorm:"size:64" json:"hash"` //去重
}

// NetPath TODO
func (ImageModel) NetPath() string {
	return fmt.Sprintf("/")
}
