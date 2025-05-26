package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/utils/pwd"
	"github.com/gin-gonic/gin"
)

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password"` //重置的密码
}

// ResetPasswordView 重置密码
func (UserApi) ResetPasswordView(c *gin.Context) {
	var cr ResetPasswordRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}

	if !global.Config.Site.Login.EmailLogin {
		resp.FailWithMsg("站点未启用邮箱注册", c)
		return
	}

	_email, _ := c.Get("email")
	email := _email.(string)

	var user model.UserModel
	err = global.DB.Take(&user, "email = ?", email).Error
	if err != nil {
		resp.FailWithMsg("不存在的用户", c)
		return
	}
	if user.RegSource != enum.RegisterEmailSourceType {
		resp.FailWithMsg("非邮箱注册用户，不能重置密码", c)
		return
	}
	hashPwd, _ := pwd.GenerateFromPassword(cr.NewPassword)
	global.DB.Model(&user).Update("password", hashPwd)
	resp.OKWithMsg("重置密码成功", c)
}
