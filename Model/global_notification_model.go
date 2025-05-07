package Model

// GlobalNotificationModel 全局通知表
type GlobalNotificationModel struct {
	Model
	Title   string `gorm:"size:64" json:"title"`
	Icon    string `gorm:"size:256" json:"icon"`
	Content string `gorm:"size:64" json:"content"`
	Href    string `gorm:"size:256" json:"href"` //消息跳转链接
}
