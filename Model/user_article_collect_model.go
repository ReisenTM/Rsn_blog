package Model

import "time"

// UserArticleCollectModel 用户个人收藏
type UserArticleCollectModel struct {
	UserID       uint         `gorm:"uniqueIndex:idx_user_article" json:"user_id"`
	UserModel    UserModel    `gorm:"foreignKey:UserID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_user_article" json:"article_id"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"-"`
	CollectID    uint         `gorm:"uniqueIndex:idx_user_article" json:"collect_id"` //属于哪一个收藏夹
	CollectModel CollectModel `gorm:"foreignKey:CollectID" json:"collect_model"`
	CreatedAt    time.Time    `json:"created_at"`
}
