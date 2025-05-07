package Model

type LogModel struct {
	Model
	Type      int8      `json:"type"`  //日志类型
	Level     int8      `json:"level"` //日志级别
	Title     string    `gorm:"size:64" json:"title"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"` //可能是访客，id设为0
	UserModel UserModel `gorm:"foreignKey:UserID" json:"-"`
	IP        string    `gorm:"size:32" json:"ip"`
	Location  string    `gorm:"size:64" json:"addr"`
	IsRead    bool      `json:"isRead"` //是否已读
}
