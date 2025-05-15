package Model

import "time"

type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IDRequest struct {
	ID uint `json:"id" form:"id" uri:"id"`
}

type DeleteRequest struct {
	IDList []uint `json:"id_list"`
}
