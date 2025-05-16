package model

import "time"

// UserTopArticleModel 置顶文章表
type UserTopArticleModel struct {
	UserID       uint         `gorm:"uniqueIndex:idx_top" json:"user_id"`
	UserModel    UserModel    `gorm:"foreignKey:UserID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_top" json:"article_id"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"-"`
	CreatedAt    time.Time    `json:"created_at"`
}
