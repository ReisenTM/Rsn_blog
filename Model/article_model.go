package Model

type ArticleModel struct {
	ArticleID    uint      `gorm:"primary_key"`
	Title        string    `gorm:"size:32" json:"title"`
	Content      string    `json:"content"`
	Preview      string    `gorm:"size:256" json:"preview"`
	CategoryID   string    `json:"category_id"`                           //分类ID
	Tags         []string  `gorm:"type:text;serializer:json" json:"tags"` //文章标签
	Cover        string    `gorm:"size:256" json:"cover"`
	UserID       uint      `json:"user_id"`
	UserModel    UserModel `gorm:"foreignKey:UserID" json:"-"`
	ViewsCount   uint      `json:"views_count"`
	FavorCount   uint      `json:"favor_count"`
	CommentCount uint      `json:"comment_count"`
	OpenComment  bool      `json:"open_comment"` //开放评论
	Status       uint      `json:"status"`       //状态:草稿，审核中，已发布
}
