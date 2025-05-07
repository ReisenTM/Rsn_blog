package Model

// UserLoginModel 用户登录表
type UserLoginModel struct {
	Model
	UserID    uint      `json:"user_id"`
	UserModel UserModel `gorm:"foreignKey:UserID" json:"-"`
	IP        string    `gorm:"size:32" json:"ip"`
	Location  string    `gorm:"size:64" json:"addr"`
	UA        string    `gorm:"size:255" json:"ua"`
}
