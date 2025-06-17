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

type SiteMsgReadRequest struct {
	ID   uint `json:"id"`
	Type int8 `json:"type"` // 一键已读的类型  1.回复评论 2.点赞收藏 3.系统
}

// SiteMsgReadView 消息已读
func (SiteMsgApi) SiteMsgReadView(c *gin.Context) {
	cr := middleware.GetBind[SiteMsgReadRequest](c)

	claims := jwts.GetClaims(c)
	if cr.ID != 0 {
		// 找这个消息是不是当前用户的
		var msg model.MessageModel
		err := global.DB.Take(&msg, "id = ? and rev_user_id = ?", cr.ID, claims.UserID).Error
		if err != nil {
			resp.FailWithMsg("消息不存在", c)
			return
		}

		if msg.IsRead {
			resp.FailWithMsg("消息已读取", c)
			return
		}

		global.DB.Model(&msg).Update("is_read", true)
		resp.OKWithMsg("消息读取成功", c)
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
	global.DB.Find(&msgList, "rev_user_id = ? and type in ? and is_read = ?", claims.UserID, typeList, false)

	if len(msgList) > 0 {
		global.DB.Model(&msgList).Update("is_read", true)
	}

	resp.OKWithMsg(fmt.Sprintf("批量读取%d条消息成功", len(msgList)), c)

}
