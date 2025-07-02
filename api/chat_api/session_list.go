package chat_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/ctype/chat_msg"
	"blogX_server/model/enum/chat_msg_type"
	"blogX_server/model/enum/relationship_enum"
	"blogX_server/service/focus_service"
	"blogX_server/utils/jwts"
	"blogX_server/utils/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type SessionListRequest struct {
	common.PageInfo
}

type SessionTable struct {
	//U1 U2 代表聊天双方
	U1        uint   `gorm:"column:U1"`
	U2        uint   `gorm:"column:U2"`
	MaxDate   string `gorm:"column:max_date"`
	Count     int    `gorm:"column:count"`
	NewChatID uint   `gorm:"column:new_chat_id"`
}

type SessionListResponse struct {
	UserID       uint                       `json:"user_id"`
	UserNickname string                     `json:"user_nickname"`
	UserAvatar   string                     `json:"user_avatar"`
	Msg          chat_msg.ChatMsg           `json:"msg"`
	MsgType      chat_msg_type.MsgType      `json:"msg_type"`
	NewMsgDate   time.Time                  `json:"new_msg_date"`
	Relation     relationship_enum.Relation `json:"relation"` // 好友关系
}

func (ChatApi) SessionListView(c *gin.Context) {
	cr := middleware.GetBind[SessionListRequest](c)
	claims := jwts.GetClaims(c)
	//先找删除的
	var deleteIDList []uint
	global.DB.Model(model.UserChatActionModel{}).Where("is_delete = ? and user_id =?", true, claims.UserID).
		Select("chat_id").Scan(&deleteIDList)
	query := global.DB.Where("")
	var column = fmt.Sprintf("(select id from chat_models where ((send_user_id = U1 and rev_user_id = U2) or (send_user_id = U2 and rev_user_id = U1))  order by created_at desc limit 1) as new_chat_id")

	if len(deleteIDList) > 0 {
		//从聊天记录排除
		query = query.Where("id not in (?)", deleteIDList)
		//从最新消息排除
		column = fmt.Sprintf("(select id from chat_models where ((send_user_id = U1 and rev_user_id = U2) or (send_user_id = U2 and rev_user_id = U1)) and id not in %s order by created_at desc limit 1) as newChatID", sql.ConvertSliceSql(deleteIDList))
	}
	var _list []SessionTable
	//查每个会话的信息
	global.DB.Model(model.ChatModel{}).
		Select(
			"least(send_user_id, rev_user_id)    as U1",
			"greatest(send_user_id, rev_user_id) as U2",
			"max(created_at)       as max_date",
			"count(*)         as c",
			column,
		).
		Where(query).
		Where("(send_user_id = ? or rev_user_id = ?)", claims.UserID, claims.UserID).
		Group("least(send_user_id, rev_user_id)").
		Group("greatest(send_user_id, rev_user_id)").
		Order("max_date desc").
		Limit(cr.GetLimit()).Offset(cr.GetOffset()).Scan(&_list)
	//查满足查询的会话个数
	var count int //聊天个数
	global.DB.Select("count(*)").Table("(?) as x",
		global.DB.
			Model(model.ChatModel{}).
			Select("count(*)").
			Where(query).
			Where("(send_user_id = ? or rev_user_id = ?)", claims.UserID, claims.UserID).
			Group("least(send_user_id, rev_user_id)").
			Group("greatest(send_user_id, rev_user_id)"),
	).Scan(&count)
	//因为返回列表包括对方用户信息，先遍历查聊天对象用户id
	var userIDList []uint
	var chatIDList []uint //显示的最新聊天 对应的id
	for _, table := range _list {
		chatIDList = append(chatIDList, table.NewChatID)
		if table.U1 == claims.UserID {
			//如果我是U1，那对方就是U2
			userIDList = append(userIDList, table.U2)
		}
		if table.U2 == claims.UserID {
			//如果我是U2，那对方就是U1
			userIDList = append(userIDList, table.U1)
		}
	}
	//建立id和内容的映射关系
	userMap := common.ScanMapV2(model.UserModel{}, common.ScanOption{
		Where: global.DB.Where("id in ?", userIDList),
	})
	chatMap := common.ScanMapV2(model.ChatModel{}, common.ScanOption{
		Where: global.DB.Where("id in ?", chatIDList),
	})
	//计算对方和用户的关系，建立关系映射
	relationMap := focus_service.CalcUserPatchRelationship(claims.UserID, userIDList)

	var list = make([]SessionListResponse, 0)
	for _, table := range _list {
		item := SessionListResponse{}
		if table.U1 == claims.UserID {
			item.UserID = table.U2
		}
		if table.U2 == claims.UserID {
			item.UserID = table.U1
		}
		item.UserNickname = userMap[item.UserID].Nickname
		item.UserAvatar = userMap[item.UserID].Avatar
		item.Msg = chatMap[table.NewChatID].Msg
		item.MsgType = chatMap[table.NewChatID].MsgType
		item.NewMsgDate = chatMap[table.NewChatID].CreatedAt
		item.Relation = relationMap[item.UserID]
		list = append(list, item)
	}
	resp.OkWithList(list, count, c)
}
