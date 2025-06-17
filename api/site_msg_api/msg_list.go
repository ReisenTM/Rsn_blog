package site_msg_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum/message_type_enum"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

type SiteMsgListRequest struct {
	common.PageInfo
	Type int8 `form:"type" binding:"required,oneof=1 2 3"` // 1评论和回复 2赞和收藏 3 系统
}
type SiteMsgListResponse struct {
	model.MessageModel
}

func (SiteMsgApi) SiteMsgListView(c *gin.Context) {
	cr := middleware.GetBind[SiteMsgListRequest](c)

	var typeList []message_type_enum.Type
	switch cr.Type {
	case 1:
		typeList = append(typeList, message_type_enum.CommentType, message_type_enum.ReplyType)
	case 2:
		typeList = append(typeList, message_type_enum.FavorArticleType, message_type_enum.FavorCommentType, message_type_enum.CollectArticleType)
	case 3:
		typeList = append(typeList, message_type_enum.SystemType)
	}

	claims := jwts.GetClaims(c)

	_list, count, _ := common.ListQuery(model.MessageModel{
		RevUserID: claims.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    global.DB.Where("type in ?", typeList),
	})

	var userIDList []uint
	for _, m := range _list {
		if m.ActionUserID != 0 {
			userIDList = append(userIDList, m.ActionUserID)
		}
	}

	var list = make([]SiteMsgListResponse, 0)
	for _, mm := range _list {
		list = append(list, SiteMsgListResponse{
			MessageModel: mm,
		})
	}

	resp.OkWithList(list, count, c)
}
