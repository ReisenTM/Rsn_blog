package chat_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/ctype/chat_msg"
	"blogX_server/model/enum/chat_msg_type"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

func (ChatApi) ChatReadView(c *gin.Context) {
	cr := middleware.GetBind[model.IDRequest](c)

	var chat model.ChatModel
	err := global.DB.Take(&chat, cr.ID).Error
	if err != nil {
		resp.FailWithMsg("消息不存在", c)
		return
	}

	item := ChatResponse{
		ChatListResponse: ChatListResponse{
			ChatModel: model.ChatModel{
				MsgType: chat_msg_type.MsgReadType,
				Msg: chat_msg.ChatMsg{
					MsgReadMsg: &chat_msg.MsgReadMsg{
						ReadChatID: chat.ID,
					},
				},
			},
		},
	}
	claims := jwts.GetClaims(c)
	var chatAc model.UserChatActionModel
	err = global.DB.Take(&chatAc, "user_id = ? and chat_id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		global.DB.Create(&model.UserChatActionModel{
			UserID: claims.UserID,
			ChatID: cr.ID,
			IsRead: true,
		})
		resp.SendWsMsg(OnlineMap, chat.SendUserID, item)
		resp.OKWithMsg("消息读取成功", c)
		return
	}

	if chatAc.IsDelete {
		resp.FailWithMsg("消息被删除", c)
		return
	}

	resp.SendWsMsg(OnlineMap, chat.SendUserID, item)
	global.DB.Model(&chatAc).Update("is_read", true)
	resp.OKWithMsg("消息读取成功", c)
}
