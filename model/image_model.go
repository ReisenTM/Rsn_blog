package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

// ImageModel 图片表
type ImageModel struct {
	Model
	Filename string `gorm:"size:64" json:"filename"`
	Path     string `gorm:"size:256" json:"path"`
	Size     int64  `json:"size"`
	Hash     string `gorm:"size:64" json:"hash"` //去重
}

// NetPath 网络请求地址拼接
func (i ImageModel) NetPath() string {
	return fmt.Sprintf("/%s", i.Path)
}

// BeforeDelete 会在删除记录前调用
func (i ImageModel) BeforeDelete(tx *gorm.DB) error {
	err := os.Remove(i.Path)
	if err != nil {
		logrus.Warnf("图片删除失败,%v", err)
	}
	//删除不了文件也不能中断,只删除记录
	return nil
}
