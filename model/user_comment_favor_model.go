package model

type UserCommentFavorModel struct {
	Model
	//联合唯一索引，一个用户只能对同一条评论点赞一次
	CommentID    uint         `gorm:"uniqueIndex:idx_name" json:"comment_id"`
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"user_id"`
	UserModel    UserModel    `gorm:"foreignKey:UserID" json:"-"`
	CommentModel CommentModel `gorm:"foreignKey:CommentID" json:"-"`
}
