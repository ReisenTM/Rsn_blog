package Model

import "time"

type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
