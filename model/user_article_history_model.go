package model

type UserArticleHistoryModel struct {
	Model
	UserID       uint         `json:"user_id"`
	ArticleID    uint         `json:"article_id"`
	UserModel    UserModel    `gorm:"foreignKey:UserID" json:"-"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"-"`
}
