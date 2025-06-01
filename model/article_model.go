package model

import (
	"blogX_server/model/ctype"
	_ "embed"
)

// ArticleModel 文章表
type ArticleModel struct {
	Model
	ArticleID    uint       `gorm:"primary_key"`
	Title        string     `gorm:"size:32" json:"title"`
	Content      string     `json:"content"`
	Preview      string     `gorm:"size:256" json:"preview"`
	CategoryID   string     `json:"category_id"`           //分类ID
	Tags         ctype.List `gorm:"type:text" json:"tags"` //文章标签
	Cover        string     `gorm:"size:256" json:"cover"`
	UserID       uint       `json:"user_id"`
	UserModel    UserModel  `gorm:"foreignKey:UserID" json:"-"`
	ViewsCount   uint       `json:"views_count"`
	FavorCount   uint       `json:"favor_count"`
	CommentCount uint       `json:"comment_count"`
	CollectCount uint       `json:"collect_count"`
	OpenComment  bool       `json:"open_comment"` //开放评论
	Status       uint       `json:"status"`       //状态:草稿，审核中，已发布
}

//go:embed mappings/article_mapping.json
var articleMapping string

func (ArticleModel) Mapping() string {
	return articleMapping
}

func (ArticleModel) Index() string {
	return "article_index"
}
