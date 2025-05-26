package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserBaseInfoResponse struct {
	ID           uint   `json:"ID"`
	CodeAge      int    `json:"codeAge"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	ViewCount    int    `json:"view_count"`
	ArticleCount int    `json:"article_count"`
	FansCount    int    `json:"fans_count"`
	FollowCount  int    `json:"follow_count"`
	Region       string `json:"region"` // ip归属地
}

func (UserApi) UserBaseInfoView(c *gin.Context) {
	var id model.IDRequest
	err := c.ShouldBindQuery(&id)
	if err != nil {
		logrus.Errorf("用户基本信息参数绑定失败%v", err)
		resp.FailWithError(err, c)
		return
	}
	var user model.UserModel
	err = global.DB.Take(&user, "id = ?", id.ID).Error
	if err != nil {
		resp.FailWithMsg("未找到用户", c)
		return
	}
	res := UserBaseInfoResponse{
		ID:           user.ID,
		CodeAge:      user.CodeAge(),
		Avatar:       user.Avatar,
		Nickname:     user.Nickname,
		ViewCount:    1,
		ArticleCount: 1, //TODO:做完文章继续
		FansCount:    0, //TODO:做完好友关系继续
		FollowCount:  0,
		Region:       user.Region,
	}

	resp.OkWithData(res, c)
}
