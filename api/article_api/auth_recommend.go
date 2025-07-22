package article_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum/relationship_enum"
	"blogX_server/service/focus_service"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthRecommendRequest struct {
	common.PageInfo
}

type AuthRecommendResponse struct {
	UserID       uint   `json:"user_id"`
	UserNickname string `json:"user_nickname"`
	UserAvatar   string `json:"user_avatar"`
	UserProfile  string `json:"user_profile"`
}

// AuthRecommendView 推荐
func (ArticleApi) AuthRecommendView(c *gin.Context) {
	cr := middleware.GetBind[AuthRecommendRequest](c)
	var userIDlist []uint
	var userCount int
	//找到所有发过文章的人
	err := global.DB.Model(model.ArticleModel{}).Group("user_id").
		Select("count(*)").Scan(&userCount).Error
	err = global.DB.Model(model.ArticleModel{}).Group("user_id").
		Offset(cr.GetOffset()).
		Limit(cr.GetLimit()).
		Select("user_id").Scan(&userIDlist).Error
	if err != nil {
		logrus.Errorf("Get user count error:%v", err)
		resp.FailWithMsg("推荐失败", c)
		return
	}
	claims, err := jwts.ParseTokenByGin(c)
	if err == nil && claims != nil {
		//登录的人就推荐陌生人和粉丝
		m := focus_service.CalcUserPatchRelationship(claims.UserID, userIDlist)
		userIDlist = []uint{}
		for u, relation := range m {
			if relation == relationship_enum.RelationStranger || relation == relationship_enum.RelationFans {
				userIDlist = append(userIDlist, u)
			}
		}
	}

	//如果没登录就推荐所有人
	var userList []model.UserModel
	global.DB.Find(&userList, "id in ?", userIDlist)
	var list = make([]AuthRecommendResponse, 0)
	for _, model := range userList {
		list = append(list, AuthRecommendResponse{
			UserID:       model.ID,
			UserNickname: model.Nickname,
			UserAvatar:   model.Avatar,
			UserProfile:  model.Profile,
		})
	}
	resp.OkWithList(list, userCount, c)
}
