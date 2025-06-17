package site_msg_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum/message_type_enum"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

type UserMsgResponse struct {
	CommentMsgCount int `json:"comment_msg_count"`
	FavorMsgCount   int `json:"favor_msg_count"`
	PrivateMsgCount int `json:"private_msg_count"`
	SystemMsgCount  int `json:"system_msg_count"`
}

// UserMsgView 用户的未读消息
func (SiteMsgApi) UserMsgView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	var msgList []model.MessageModel
	global.DB.Find(&msgList, "rev_user_id = ? and is_read = ?", claims.UserID, false)

	//统计未读数量
	var data UserMsgResponse
	for _, model := range msgList {
		switch model.Type {
		case message_type_enum.CommentType, message_type_enum.ReplyType:
			data.CommentMsgCount++
		case message_type_enum.FavorArticleType, message_type_enum.FavorCommentType, message_type_enum.CollectArticleType:
			data.FavorMsgCount++
		case message_type_enum.SystemType:
			data.SystemMsgCount++
		}
	}
	//TODO:找未读的私信
	//var chatList []models.ChatModel
	//// 接收人是我，而且这个消息未读
	//global.DB.Find(&chatList, "rev_user_id = ?", claims.UserID)
	//var chatIDList []uint
	//for _, model := range chatList {
	//	chatIDList = append(chatIDList, model.ID)
	//}
	//chatAcMap := common.ScanMapV2(model.UserChatActionModel{}, common.ScanOption{
	//	Where: global.DB.Where("chat_id in ?", chatIDList),
	//})
	//for _, model := range chatList {
	//	_, ok := chatAcMap[model.ID]
	//	if !ok {
	//		data.PrivateMsgCount++
	//		continue
	//	}
	//}
	// 过滤掉删除的，只取未读的
	var userReadMsgIDList []uint
	global.DB.Model(model.UserGlobalNotificationModel{}).
		Where("user_id = ? and (is_read = ? or is_delete = ?)", claims.UserID, true, true).
		Select("id").Scan(&userReadMsgIDList)
	// 算未读的全局消息
	var systemMsg []model.GlobalNotificationModel
	query := global.DB.Where("")
	if len(userReadMsgIDList) > 0 {
		query.Where("id not in ?", userReadMsgIDList)
	}
	global.DB.Where(query).Find(&systemMsg)
	data.SystemMsgCount += len(systemMsg)

	resp.OkWithData(data, c)
}
