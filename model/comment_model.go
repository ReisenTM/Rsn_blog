package model

type CommentModel struct {
	Model
	Content      string          `gorm:"size:256" json:"content"`
	UserID       uint            `json:"user_id"`
	ArticleID    uint            `json:"article_id"`
	UserModel    UserModel       `gorm:"foreignKey:UserID" json:"-"`
	ArticleModel ArticleModel    `gorm:"foreignKey:ArticleID" json:"-"`
	ParentID     *uint           `json:"parent_id"` //父评论id
	ParentModel  *CommentModel   `gorm:"foreignKey:ParentID" json:"-"`
	SubComment   []*CommentModel `gorm:"foreignKey:ParentID" json:"-"` //子评论列表
	RootID       *uint           `json:"root_id"`                      //根评论id
	FavorCount   int             `json:"favor_count"`                  //点赞统计
	ReplyCount   int             `json:"reply_count"`                  //回复统计
}
