package user_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"github.com/gin-gonic/gin"
	"time"
)

type UserListRequest struct {
	common.PageInfo
}
type UserListResponse struct {
	Username       string        `json:"username"`
	ID             uint          `json:"id"`
	IP             string        `json:"ip"`
	Region         string        `json:"addr"`
	Nickname       string        `json:"nickname"`
	Avatar         string        `json:"avatar"`
	ArticleCount   int           `json:"article_count"`
	FansCount      int           `json:"fans_count"`
	FollowersCount int           `json:"followers_count"`
	IndexCount     int           `json:"index_count"` //主页访问数
	CreatedAt      time.Time     `json:"created_at"`  //创建时间
	LastLogin      *time.Time    `json:"last_login"`  //上次登录时间
	Role           enum.RoleType `json:"role"`
}

func (UserApi) UserListView(c *gin.Context) {
	cr := middleware.GetBind[UserListRequest](c)

	_list, count, _ := common.ListQuery(model.UserModel{}, common.Options{
		Likes:    []string{"nickname", "username"},
		Preloads: []string{"ArticleList", "LoginList"},
		PageInfo: cr.PageInfo,
	})
	var list = make([]UserListResponse, 0)
	for _, model := range _list {
		item := UserListResponse{
			ID:           model.ID,
			Nickname:     model.Nickname,
			Username:     model.Username,
			Avatar:       model.Avatar,
			IP:           model.IP,
			Region:       model.Region,
			ArticleCount: len(model.ArticleList),
			IndexCount:   1000,
			CreatedAt:    model.CreatedAt,
			Role:         model.Role,
		}
		if len(model.LoginList) > 0 {
			item.LastLogin = &model.LoginList[len(model.LoginList)-1].CreatedAt
		}
		list = append(list, item)
	}

	resp.OkWithList(list, count, c)
}
