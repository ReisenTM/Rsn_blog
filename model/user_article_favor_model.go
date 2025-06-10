package model

import "time"

// UserArticleFavorModel 文章点赞表
type UserArticleFavorModel struct {
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"user_id"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_name" json:"article_id"`
	CreatedAt    time.Time    `json:"created_at"`                    //点赞时间
	UserModel    UserModel    `gorm:"foreignKey:UserID" json:"-"`    //关联的用户
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"-"` //关联的文章
}
