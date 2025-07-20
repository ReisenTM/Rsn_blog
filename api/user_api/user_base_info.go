package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum/relationship_enum"
	"blogX_server/service/focus_service"
	"blogX_server/service/redis_service/redis_user"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

type UserBaseInfoResponse struct {
	UserID       uint                       `json:"user_id"`
	CodeAge      int                        `json:"code_age"`
	Avatar       string                     `json:"avatar"`
	Nickname     string                     `json:"nickname"`
	ViewsCount   int                        `json:"views_count"`
	ArticleCount int                        `json:"article_count"`
	FansCount    int                        `json:"fans_count"`
	FollowsCount int                        `json:"follows_count"`
	Region       string                     `json:"region"`        // ip归属地
	OpenCollect  bool                       `json:"open_collect"`  // 公开我的收藏
	OpenFollows  bool                       `json:"open_follows"`  // 公开我的关注
	OpenFans     bool                       `json:"open_fans"`     // 公开我的粉丝
	HomeStyleID  uint                       `json:"home_style_id"` // 主页样式的id
	Relation     relationship_enum.Relation `json:"relation"`      // 与登录人的关系
}

func (UserApi) UserBaseInfoView(c *gin.Context) {
	var cr model.IDRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}

	var user model.UserModel
	err = global.DB.Preload("UserConfModel").Preload("ArticleList").Take(&user, cr.ID).Error
	if err != nil {
		resp.FailWithMsg("不存在的用户", c)
		return
	}
	data := UserBaseInfoResponse{
		UserID:       user.ID,
		CodeAge:      user.CodeAge(),
		Avatar:       user.Avatar,
		Nickname:     user.Nickname,
		ViewsCount:   user.UserConfModel.ViewCount + redis_user.GetCacheLook(cr.ID),
		ArticleCount: len(user.ArticleList),
		FansCount:    0,
		FollowsCount: 0,
		Region:       user.Region,
		OpenCollect:  user.UserConfModel.OpenCollection,
		OpenFollows:  user.UserConfModel.OpenFollows,
		OpenFans:     user.UserConfModel.OpenFans,
		HomeStyleID:  user.UserConfModel.HomeStyle,
	}
	//计算关系
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		data.Relation = focus_service.CalcUserRelationship(claims.UserID, cr.ID)
	}
	//计算粉丝和关注数
	var focusList []model.UserFocusModel
	global.DB.Find(&focusList, "user_id = ? or focus_user_id = ?", cr.ID, cr.ID)
	for _, model := range focusList {
		if model.UserID == cr.ID {
			data.FansCount++
		}
		if model.FocusUserID == cr.ID {
			data.FollowsCount++
		}
	}

	redis_user.SetCacheLook(cr.ID, true)

	resp.OkWithData(data, c)
}
