package model

import "time"

type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IDRequest 用于指定请求(文章）
type IDRequest struct {
	ID uint `json:"id" form:"id" uri:"id"`
}

// RemoveRequest 用于批量删除请求
type RemoveRequest struct {
	IDList []uint `json:"id_list"`
}
