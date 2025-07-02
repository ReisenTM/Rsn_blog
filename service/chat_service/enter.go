package chat_service

import (
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/ctype/chat_msg"
	"blogX_server/model/enum/chat_msg_type"
	"blogX_server/utils/xss"
	"github.com/sirupsen/logrus"
)

// ToChat A给B发消息
func ToChat(A, B uint, msgType chat_msg_type.MsgType, msg chat_msg.ChatMsg) {
	err := global.DB.Create(&model.ChatModel{
		SendUserID: A,
		RevUserID:  B,
		MsgType:    msgType,
		Msg:        msg,
	}).Error
	if err != nil {
		logrus.Errorf("对话创建失败 %s", err)
	}
}

func ToTextChat(A, B uint, content string) {
	ToChat(A, B, chat_msg_type.TextMsgType, chat_msg.ChatMsg{
		TextMsg: &chat_msg.TextMsg{
			Content: content,
		},
	})
}

func ToImageChat(A, B uint, src string) {
	ToChat(A, B, chat_msg_type.ImageMsgType, chat_msg.ChatMsg{
		ImageMsg: &chat_msg.ImageMsg{
			Src: src,
		},
	})
}

func ToMarkdownChat(A, B uint, content string) {
	// 过滤xss
	filterContent := xss.XssFilter(content)
	ToChat(A, B, chat_msg_type.MarkdownMsgType, chat_msg.ChatMsg{
		MarkdownMsg: &chat_msg.MarkdownMsg{
			Content: filterContent,
		},
	})
}
