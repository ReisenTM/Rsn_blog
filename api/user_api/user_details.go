package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
	"time"
)

type UserDetailResponse struct {
	ID             uint                    `json:"id"`
	CreatedAt      time.Time               `json:"createdAt"`
	Username       string                  `json:"username"`
	Nickname       string                  `json:"nickname"`
	Avatar         string                  `json:"avatar"`
	Profile        string                  `json:"profile"`
	RegisterSource enum.RegisterSourceType `json:"register_source"` // 注册来源
	CodeAge        int                     `json:"codeAge"`         // 码龄
	Role           enum.RoleType           `json:"role"`            // 角色
	model.UserConfModel
	Email       string `json:"email"`
	UsePassword bool   `json:"usePassword"`
}

// UserDetailView 用户详情
func (UserApi) UserDetailView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	var user model.UserModel
	err := global.DB.Preload("UserConfModel").Take(&user, claims.UserID).Error
	if err != nil {
		resp.FailWithMsg("用户不存在", c)
		return
	}

	var data = UserDetailResponse{
		ID:             user.ID,
		CreatedAt:      user.CreatedAt,
		Username:       user.Username,
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Profile:        user.Profile,
		Role:           user.Role,
		RegisterSource: user.RegSource,
		CodeAge:        user.CodeAge(),
		Email:          user.Email,
	}
	if user.Password != "" {
		data.UsePassword = true
	}
	if user.UserConfModel != nil {
		data.UserConfModel = *user.UserConfModel
	}

	resp.OkWithData(data, c)
}
