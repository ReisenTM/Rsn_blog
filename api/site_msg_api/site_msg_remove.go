package site_msg_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum/message_type_enum"
	"blogX_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SiteMsgRemoveRequest struct {
	ID   uint `json:"id"`
	Type int8 `json:"type"` //一键删除的类型
}

// SiteMsgRemoveView 删除消息
func (SiteMsgApi) SiteMsgRemoveView(c *gin.Context) {
	cr := middleware.GetBind[SiteMsgRemoveRequest](c)
	claims := jwts.GetClaims(c)
	var msg model.MessageModel

	//看消息是否是当前用户的
	if cr.ID != 0 {
		err := global.DB.Take(&msg, "id = ? and rev_user_id = ?", cr.ID, claims.UserID).Error
		if err != nil {
			resp.FailWithMsg("消息不存在", c)
			return
		}
		global.DB.Delete(&msg)
		resp.OKWithMsg("消息删除成功", c)
		return
	}
	var typeList []message_type_enum.Type
	switch cr.Type {
	case 1:
		typeList = append(typeList, message_type_enum.CommentType, message_type_enum.ReplyType)
	case 2:
		typeList = append(typeList, message_type_enum.FavorArticleType, message_type_enum.FavorCommentType, message_type_enum.CollectArticleType)
	case 3:
		typeList = append(typeList, message_type_enum.SystemType)
	}
	var msgList []model.MessageModel
	global.DB.Find(&msgList, "rev_user_id = ? and type in ?", claims.UserID, typeList)
	if len(msgList) > 0 {
		global.DB.Delete(&msgList)
	}
	resp.OKWithMsg(fmt.Sprintf("批量删除%d条消息成功", len(msgList)), c)
}
