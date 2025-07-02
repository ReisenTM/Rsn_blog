package model

import (
	"blogX_server/model/ctype/chat_msg"
	"blogX_server/model/enum/chat_msg_type"
)

type ChatModel struct {
	Model
	SendUserID    uint                  `json:"send_user_id"`
	SendUserModel UserModel             `gorm:"foreignKey:SendUserID"  json:"-"`
	RevUserID     uint                  `json:"rev_user_id"`
	RevUserModel  UserModel             `gorm:"foreignKey:RevUserID"  json:"-"`
	MsgType       chat_msg_type.MsgType `json:"msg_type"` // 消息类型
	Msg           chat_msg.ChatMsg      `gorm:"type:longtext;serializer:json" json:"msg"`
}
