package focus_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum/relationship_enum"
	"blogX_server/service/focus_service"
	"blogX_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type FocusApi struct {
}
type FocusUserRequest struct {
	FocusUserID uint `json:"focus_user_id" binding:"required"`
}

// FocusUserView 登录人关注用户
func (FocusApi) FocusUserView(c *gin.Context) {
	cr := middleware.GetBind[FocusUserRequest](c)

	claims := jwts.GetClaims(c)
	if cr.FocusUserID == claims.UserID {
		resp.FailWithMsg("你时刻都在关注自己", c)
		return
	}
	// 查关注的用户是否存在
	var user model.UserModel
	err := global.DB.Take(&user, cr.FocusUserID).Error
	if err != nil {
		resp.FailWithMsg("关注用户不存在", c)
		return
	}

	// 查之前是否已经关注过他了
	var focus model.UserFocusModel
	err = global.DB.Take(&focus, "user_id = ? and focus_user_id = ?", claims.UserID, user.ID).Error
	if err == nil {
		resp.FailWithMsg("请勿重复关注", c)
		return
	}

	// 每天关注是不是应该有个限度？
	// 每天的取关也要有个限度？

	// 关注
	global.DB.Create(&model.UserFocusModel{
		UserID:      claims.UserID,
		FocusUserID: cr.FocusUserID,
	})

	resp.OKWithMsg("关注成功", c)
	return
}

type FocusUserListRequest struct {
	common.PageInfo
	FocusUserID uint `form:"focus_user_id"`
	UserID      uint `form:"userID"` // 查用户的关注
}
type UserListResponse struct {
	UserID       uint                       `json:"user_id"`
	UserNickname string                     `json:"user_nickname"`
	UserAvatar   string                     `json:"user_avatar"`
	UserProfile  string                     `json:"user_profile"`
	Relationship relationship_enum.Relation `json:"relationship"`
	CreatedAt    time.Time                  `json:"createdAt"`
}

// FocusUserListView 我的关注和用户的关注
func (FocusApi) FocusUserListView(c *gin.Context) {
	cr := middleware.GetBind[FocusUserListRequest](c)

	claims, err := jwts.ParseTokenByGin(c)

	if cr.UserID != 0 {
		// 传了用户id，我就查这个人关注的用户列表
		var userConf model.UserConfModel
		err1 := global.DB.Take(&userConf, "user_id = ?", cr.UserID).Error
		if err1 != nil {
			resp.FailWithMsg("用户配置信息不存在", c)
			return
		}
		if !userConf.OpenFollows {
			resp.FailWithMsg("此用户未公开我的关注", c)
			return
		}

		// 如果你没登录。我就不允许你查第二页
		if err != nil || claims == nil {
			if cr.Limit > 10 || cr.Page > 1 {
				resp.FailWithMsg("未登录用户只能显示第一页", c)
				return
			}
		}

	} else {
		if err != nil || claims == nil {
			resp.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.UserID
	}

	query := global.DB.Where("")
	if cr.Key != "" {
		// 模糊匹配用户
		var userIDList []uint
		global.DB.Model(&model.UserModel{}).
			Where("nickname like ?", fmt.Sprintf("%%%s%%", cr.Key)).
			Select("id").Scan(&userIDList)
		if len(userIDList) > 0 {
			query.Where("focus_user_id in ?", userIDList)
		}
	}

	_list, count, _ := common.ListQuery(model.UserFocusModel{
		FocusUserID: cr.FocusUserID,
		UserID:      cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    query,
		Preloads: []string{"FocusUserModel"},
	})

	var m = map[uint]relationship_enum.Relation{}
	if err == nil && claims != nil {
		var userIDList []uint
		for _, i2 := range _list {
			userIDList = append(userIDList, i2.FocusUserID)
		}
		m = focus_service.CalcUserPatchRelationship(claims.UserID, userIDList)

	}

	var list = make([]UserListResponse, 0)
	for _, model := range _list {
		list = append(list, UserListResponse{
			UserID:       model.FocusUserID,
			UserNickname: model.FocusUserModel.Nickname,
			UserAvatar:   model.FocusUserModel.Avatar,
			UserProfile:  model.FocusUserModel.Profile,
			Relationship: m[model.FocusUserID],
			CreatedAt:    model.CreatedAt,
		})
	}

	resp.OkWithList(list, count, c)
}

// FansUserListView 我的粉丝和用户的粉丝
func (FocusApi) FansUserListView(c *gin.Context) {
	cr := middleware.GetBind[FocusUserListRequest](c)
	claims, err := jwts.ParseTokenByGin(c)
	if cr.UserID != 0 {
		// 传了用户id，我就查这个人的粉丝列表
		var userConf model.UserConfModel
		err1 := global.DB.Take(&userConf, "user_id = ?", cr.UserID).Error
		if err1 != nil {
			resp.FailWithMsg("用户配置信息不存在", c)
			return
		}
		if !userConf.OpenFans {
			resp.FailWithMsg("此用户未公开我的粉丝", c)
			return
		}
		// 如果你没登录。我就不允许你查第二页
		if err != nil || claims == nil {
			if cr.Limit > 10 || cr.Page > 1 {
				resp.FailWithMsg("未登录用户只能显示第一页", c)
				return
			}
		}
	} else {
		if err != nil || claims == nil {
			resp.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.UserID
	}

	query := global.DB.Where("")
	if cr.Key != "" {
		// 模糊匹配用户
		var userIDList []uint
		global.DB.Model(&model.UserModel{}).
			Where("nickname like ?", fmt.Sprintf("%%%s%%", cr.Key)).
			Select("id").Scan(&userIDList)
		if len(userIDList) > 0 {
			query.Where("user_id in ?", userIDList)
		}
	}

	_list, count, _ := common.ListQuery(model.UserFocusModel{
		FocusUserID: cr.UserID,
		UserID:      cr.FocusUserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    query,
		Preloads: []string{"UserModel"},
	})
	var m = map[uint]relationship_enum.Relation{}
	if err == nil && claims != nil {
		var userIDList []uint
		for _, i2 := range _list {
			userIDList = append(userIDList, i2.UserID)
		}
		m = focus_service.CalcUserPatchRelationship(claims.UserID, userIDList)
	}

	var list = make([]UserListResponse, 0)
	for _, model := range _list {
		list = append(list, UserListResponse{
			UserID:       model.UserID,
			UserNickname: model.UserModel.Nickname,
			UserAvatar:   model.UserModel.Avatar,
			UserProfile:  model.UserModel.Profile,
			CreatedAt:    model.CreatedAt,
			Relationship: m[model.UserID],
		})
	}

	resp.OkWithList(list, count, c)
}

// UnFocusUserView 登录人取关用户
func (FocusApi) UnFocusUserView(c *gin.Context) {
	cr := middleware.GetBind[FocusUserRequest](c)

	claims := jwts.GetClaims(c)
	if cr.FocusUserID == claims.UserID {
		resp.FailWithMsg("你无法取关自己", c)
		return
	}
	// 查关注的用户是否存在
	var user model.UserModel
	err := global.DB.Take(&user, cr.FocusUserID).Error
	if err != nil {
		resp.FailWithMsg("取关用户不存在", c)
		return
	}

	// 查之前是否已经关注过他了
	var focus model.UserFocusModel
	err = global.DB.Take(&focus, "user_id = ? and focus_user_id = ?", claims.UserID, user.ID).Error
	if err != nil {
		resp.FailWithMsg("未关注此用户", c)
		return
	}
	// 每天的取关也要有个限度？
	// 取关
	global.DB.Delete(&focus)
	resp.OKWithMsg("取消关注成功", c)
	return
}
