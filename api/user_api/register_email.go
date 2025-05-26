package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/service/user_service"
	"blogX_server/utils/jwts"
	"blogX_server/utils/pwd"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
)

type RegisterEmailRequest struct {
	EmailID   string `json:"email_id" binding:"required"`   //冗余
	EmailCode string `json:"email_code" binding:"required"` //冗余
	Password  string `json:"password"`
}

func (UserApi) RegisterEmailView(c *gin.Context) {
	var req RegisterEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("参数绑定失败,%s", err.Error())
		resp.FailWithError(err, c)
		return
	}
	if !global.Config.Site.Login.EmailLogin {
		resp.FailWithMsg("站点未启用邮箱注册", c)
		return
	}

	// 创建用户
	uname := base64Captcha.RandText(5, "0123456789")
	hashPwd, _ := pwd.GenerateFromPassword(req.Password)
	_email, _ := c.Get("email")
	email := _email.(string)
	var user = model.UserModel{
		Username:  fmt.Sprintf("b_%s", uname),
		Nickname:  "邮箱登录用户",
		Email:     email,
		Password:  hashPwd,
		RegSource: enum.RegisterEmailSourceType,
		Role:      enum.RoleUserType,
	}
	err := global.DB.Create(&user).Error
	if err != nil {
		resp.FailWithMsg("邮箱注册失败", c)
		return
	}
	// 颁发token
	token, err := jwts.GetToken(jwts.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	})
	if err != nil {
		resp.FailWithMsg("邮箱登录失败", c)
		return
	}
	user_service.NewUserService(user).UserLogin(c)
	resp.OkWithData(token, c)
}
