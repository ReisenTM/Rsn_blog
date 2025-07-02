package chat_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

// UserChatDeleteView 删除选中的的聊天记录
func (ChatApi) UserChatDeleteView(c *gin.Context) {
	cr := middleware.GetBind[model.RemoveRequest](c)
	//先查对应的聊天记录
	var chatList []model.ChatModel
	global.DB.Find(&chatList, "id in ?", cr.IDList)
	claims := jwts.GetClaims(c)
	//
	chatMap := common.ScanMapV2(model.UserChatActionModel{}, common.ScanOption{
		Where: global.DB.Where("user_id = ? and chat_id in ?", claims.UserID, cr.IDList),
		Key:   "ChatID",
	})
	var addChatAc []model.UserChatActionModel
	var updateChatAcIdList []uint
	for _, mod := range chatList {
		// 判断这个消息在不在行为表里
		chat, ok := chatMap[mod.ID]
		if !ok {
			// 找不到的情况
			addChatAc = append(addChatAc, model.UserChatActionModel{
				UserID:   claims.UserID,
				ChatID:   mod.ID,
				IsDelete: true,
			})
			continue
		}
		if chat.IsDelete {
			//如果在且已经是删除状态
			continue
		}
		updateChatAcIdList = append(updateChatAcIdList, chat.ID)
	}

	if len(addChatAc) > 0 {
		err := global.DB.Create(&addChatAc).Error
		if err != nil {
			resp.FailWithMsg("删除消息失败", c)
			return
		}
	}
	if len(updateChatAcIdList) > 0 {
		//更新状态
		global.DB.Model(&model.UserChatActionModel{}).Where("id in ?", updateChatAcIdList).Update("is_delete", true)
	}
	resp.OKWithMsg("消息删除成功", c)

}
