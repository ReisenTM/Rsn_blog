package model

type UserFocusModel struct {
	Model
	UserID         uint      `json:"user_id"` // 用户id
	UserModel      UserModel `gorm:"foreignKey:UserID" json:"-"`
	FocusUserID    uint      `json:"focus_user_id"` // 关注的用户
	FocusUserModel UserModel `gorm:"foreignKey:FocusUserID" json:"-"`
}
