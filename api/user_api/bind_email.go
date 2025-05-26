package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

// BindEmailView 绑定邮箱
func (UserApi) BindEmailView(c *gin.Context) {

	if !global.Config.Site.Login.EmailLogin {
		resp.FailWithMsg("站点未启用邮箱注册", c)
		return
	}

	_email, _ := c.Get("email")
	email := _email.(string)

	user, err := jwts.GetClaims(c).GetUser()
	if err != nil {
		resp.FailWithMsg("不存在的用户", c)
		return
	}
	global.DB.Model(&user).Update("email", email)
	resp.OKWithMsg("邮箱绑定成功", c)
}
