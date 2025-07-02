package chat_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

type ChatApi struct {
}

type ChatListRequest struct {
	common.PageInfo
	SendUserID uint `form:"sendUserID"`
	RevUserID  uint `form:"revUserID"  binding:"required"`
	Type       int8 `form:"type" binding:"required,oneof=1 2"` //用户调 管理员调
}

type ChatListResponse struct {
	model.ChatModel
	SendUserNickname string `json:"send_user_nickname"`
	SendUserAvatar   string `json:"send_user_avatar"`
	RevUserNickname  string `json:"rev_user_nickname"`
	RevUserAvatar    string `json:"rev_user_avatar"`
	IsMe             bool   `json:"is_me"`   // 是我发的
	IsRead           bool   `json:"is_read"` // 消息是否已读
}

func (ChatApi) ChatListView(c *gin.Context) {
	cr := middleware.GetBind[ChatListRequest](c)

	claims := jwts.GetClaims(c)
	var deletedIDList []uint
	var userChatActionList []model.UserChatActionModel
	var chatReadMap = map[uint]bool{}
	//找所有未删除的
	global.DB.Find(&userChatActionList, "user_id = ? and (is_delete = ? or is_delete is null)", cr.RevUserID, 0)
	for _, model := range userChatActionList {
		chatReadMap[model.ChatID] = true
	}

	switch cr.Type {
	case 1: // 前台用户调的
		cr.SendUserID = claims.UserID
		// 找我删除的消息
		global.DB.Model(model.UserChatActionModel{}).
			Where("user_id = ? and is_delete = ?", claims.UserID, true).
			Select("chat_id").Scan(&deletedIDList)
	case 2:
		if claims.Role != enum.RoleAdminType {
			resp.FailWithMsg("权限错误", c)
			return
		}
		if cr.SendUserID == 0 {
			resp.FailWithMsg("sendUserID必填", c)
			return
		}
	}
	//找双方收发信息
	query := global.DB.Where("(send_user_id = ? and rev_user_id = ?) or(send_user_id = ? and rev_user_id = ?) ",
		cr.SendUserID, cr.RevUserID, cr.RevUserID, cr.SendUserID,
	)
	//剔除删除的
	if len(deletedIDList) > 0 {
		query.Where("id not in ?", deletedIDList)
	}
	//时间降序排列
	cr.Order = "created_at desc"
	_list, count, _ := common.ListQuery(model.ChatModel{}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"SendUserModel", "RevUserModel"},
		Where:    query,
	})

	var list = make([]ChatListResponse, 0)
	for _, model := range _list {
		item := ChatListResponse{
			ChatModel:        model,
			SendUserNickname: model.SendUserModel.Nickname,
			SendUserAvatar:   model.SendUserModel.Avatar,
			RevUserNickname:  model.RevUserModel.Nickname,
			RevUserAvatar:    model.RevUserModel.Nickname,
			IsRead:           chatReadMap[model.ID],
		}
		//标记哪些是我发出的
		if model.SendUserID == claims.UserID {
			item.IsMe = true
		}
		list = append(list, item)
	}

	resp.OkWithList(list, count, c)
}
