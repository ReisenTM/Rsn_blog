package model

// UserChatActionModel 用户聊天action：已读，删除
type UserChatActionModel struct {
	Model
	UserID   uint `json:"user_id"`
	ChatID   uint `json:"chat_id"`
	IsRead   bool `json:"is_read"`
	IsDelete bool `json:"is_delete"`
}
