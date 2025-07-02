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

// UserChatDeleteByUserView 删除和目标的全部聊天记录
func (ChatApi) UserChatDeleteByUserView(c *gin.Context) {
	cr := middleware.GetBind[model.IDRequest](c)

	var user model.UserModel
	err := global.DB.Take(&user, cr.ID).Error
	if err != nil {
		resp.FailWithMsg("用户不存在", c)
		return
	}
	claims := jwts.GetClaims(c)

	// 找我和他产生了哪些消息
	query := global.DB.Where("(send_user_id = ? and rev_user_id = ?) or(send_user_id = ? and rev_user_id = ?) ",
		claims.UserID, cr.ID, cr.ID, claims.UserID,
	)

	var chatList []model.ChatModel
	global.DB.Where(query).Find(&chatList)

	var idList []uint
	for _, mod := range chatList {
		idList = append(idList, mod.ID)
	}

	chatMap := common.ScanMapV2(model.UserChatActionModel{}, common.ScanOption{
		Where: global.DB.Where("user_id = ? and chat_id in ?", claims.UserID, idList),
		Key:   "ChatID",
	})

	var addChatAc []model.UserChatActionModel
	var updateChatAcIdList []uint
	for _, mod := range chatList {
		// 判断这个消息是不是删过了
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
		global.DB.Model(&model.UserChatActionModel{}).Where("id in ?", updateChatAcIdList).Update("is_delete", true)
	}
	resp.OKWithMsg("消息删除成功", c)
}
