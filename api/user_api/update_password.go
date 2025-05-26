package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model/enum"
	"blogX_server/utils/jwts"
	"blogX_server/utils/pwd"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UpdateUserPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UpdatePasswordView 更新密码
func (UserApi) UpdatePasswordView(c *gin.Context) {
	var cr UpdateUserPasswordRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		resp.FailWithError(err, c)
		logrus.Errorf("UpdatePassword请求参数绑定失败:%v", err)
		return
	}
	cl := jwts.GetClaims(c)
	user, err := cl.GetUser()
	if err != nil {
		resp.FailWithMsg("用户不存在", c)
		return
	}
	// 邮箱注册的、绑了邮箱的
	if !(user.RegSource == enum.RegisterEmailSourceType || user.Email != "") {
		resp.FailWithMsg("仅支持邮箱注册或绑定邮箱的用户修改密码", c)
		return
	}

	// 校验之前的密码
	if !pwd.CompareHashAndPassword(user.Password, cr.OldPassword) {
		resp.FailWithMsg("旧密码错误", c)
		return
	}

	hashPwd, _ := pwd.GenerateFromPassword(cr.NewPassword)
	global.DB.Model(&user).Update("password", hashPwd)
	resp.OKWithMsg("修改密码成功", c)
}
