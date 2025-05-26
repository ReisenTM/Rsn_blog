package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/service/user_service"
	"blogX_server/utils/jwts"
	"blogX_server/utils/pwd"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PwdLoginRequest struct {
	Account  string `json:"account"` //可以是邮箱或者用户名
	Password string `json:"password"`
}

// PwdLoginView 账号密码登录
func (UserApi) PwdLoginView(c *gin.Context) {
	var cr PwdLoginRequest
	err := c.ShouldBind(&cr)
	if err != nil {
		resp.FailWithMsg("用户账号密码绑定失败", c)
		logrus.Errorf("用户账号密码绑定失败,%v", err)
		return
	}
	if !global.Config.Site.Login.UsernamePwdLogin {
		resp.FailWithMsg("站点未启用账号密码登录", c)
		return
	}
	var user model.UserModel
	err = global.DB.Take(&user, "(username = ? or email = ?)and password <> ''", cr.Account, cr.Account).Error
	if err != nil {
		resp.FailWithMsg("账号不存在，请注册", c)
		return
	}
	ok := pwd.CompareHashAndPassword(user.Password, cr.Password)
	if !ok {
		resp.FailWithMsg("密码错误，请重新输入", c)
		return
	}
	//颁发token
	token, err := jwts.GetToken(jwts.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
	user_service.NewUserService(user).UserLogin(c)

	resp.OkWithData(token, c)
}
