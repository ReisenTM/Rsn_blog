package model

// CategoryModel 分类表
type CategoryModel struct {
	Model
	Title     string    `gorm:"size:32" json:"title"`
	UserID    uint      `json:"user_id"`
	UserModel UserModel `gorm:"foreignKey:UserID" json:"-"`
}
