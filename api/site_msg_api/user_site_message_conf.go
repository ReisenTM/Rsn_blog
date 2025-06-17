package site_msg_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/utils/jwts"
	"blogX_server/utils/mps"
	"github.com/gin-gonic/gin"
)

// UserSiteMessageConfView 用户消息配置查看
func (SiteMsgApi) UserSiteMessageConfView(c *gin.Context) {
	claims := jwts.GetClaims(c)

	var userMsgConf model.UserMessageConfModel
	err := global.DB.Take(&userMsgConf, "user_id = ?", claims.UserID).Error
	if err != nil {
		resp.FailWithMsg("用户消息配置不存在", c)
		return
	}

	resp.OkWithData(userMsgConf, c)
}

type UserMessageConfUpdateRequest struct {
	OpenCommentMessage *bool `json:"open_comment_message" u:"open_comment_message"` // 开启回复和评论
	OpenFavorMessage   *bool `json:"open_favor_message" u:"open_favor_message"`     // 开启赞和收藏
	OpenPrivateChat    *bool `json:"open_private_chat" u:"open_private_chat"`       // 是否开启私聊
}

// UserSiteMessageConfUpdateView 用户站内信配置更新
func (SiteMsgApi) UserSiteMessageConfUpdateView(c *gin.Context) {
	var cr = middleware.GetBind[UserMessageConfUpdateRequest](c)
	claims := jwts.GetClaims(c)

	var userMsgConf model.UserMessageConfModel
	err := global.DB.Take(&userMsgConf, "user_id = ?", claims.UserID).Error
	if err != nil {
		resp.FailWithMsg("用户消息配置不存在", c)
		return
	}

	mp := mps.StructToMap(cr, "u")
	global.DB.Model(&userMsgConf).Updates(mp)

	resp.OKWithMsg("用户消息配置更新成功", c)
}
