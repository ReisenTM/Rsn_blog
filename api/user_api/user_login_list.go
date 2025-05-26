package user_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
	"time"
)

type UserLoginListRequest struct {
	common.PageInfo
	UserID    uint   `form:"userId"`
	Ip        string `form:"ip"`
	Location  string `form:"location"`
	StartTime string `form:"startTime"` // 起止时间的 年月日时分秒格式
	EndTime   string `form:"endTime"`
	Type      int8   `form:"type" binding:"required,oneof=1 2"` // 1 用户：只能查自己的  2 管理员 ：能查全部
}
type UserLoginListResponse struct {
	model.UserLoginModel
	UserNickname string `json:"user_nickname,omitempty"`
	UserAvatar   string `json:"user_avatar,omitempty"`
}

// UserLoginListView 登录列表查询
func (UserApi) UserLoginListView(c *gin.Context) {
	var cr UserLoginListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}

	claims := jwts.GetClaims(c)
	if cr.Type == 1 {
		cr.UserID = claims.UserID
	}

	var query = global.DB.Where("")
	if cr.StartTime != "" {
		_, err = time.Parse("2006-01-02 15:04:05", cr.StartTime)
		if err != nil {
			resp.FailWithMsg("开始时间格式错误", c)
			return
		}
		query.Where("created_at >= ?", cr.StartTime)
	}
	if cr.EndTime != "" {
		_, err = time.Parse("2006-01-02 15:04:05", cr.EndTime)
		if err != nil {
			resp.FailWithMsg("结束时间格式错误", c)
			return
		}
		query.Where("created_at <= ?", cr.EndTime)
	}
	var preloads []string
	if cr.Type == 2 {
		preloads = []string{"UserModel"}
	}

	_list, count, _ := common.ListQuery(model.UserLoginModel{
		UserID:   cr.UserID,
		IP:       cr.Ip,
		Location: cr.Location,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    query,
		Preloads: preloads,
	})
	//结果列表
	var list = make([]UserLoginListResponse, 0)
	for _, m := range _list {
		list = append(list, UserLoginListResponse{
			UserLoginModel: m,
			UserNickname:   m.UserModel.Nickname,
			UserAvatar:     m.UserModel.Avatar,
		})
	}

	resp.OkWithList(list, count, c)

}
