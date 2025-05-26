package user_service

import (
	"blogX_server/core"
	"blogX_server/global"
	"blogX_server/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserLogin 用户登录日志
func (u *UserService) UserLogin(c *gin.Context) {
	ip := c.ClientIP()
	location := core.GetIPLoc(ip)
	ua := c.GetHeader("User-Agent")
	err := global.DB.Create(&model.UserLoginModel{
		UserID:   u.userModel.ID,
		IP:       ip,
		Location: location,
		UA:       ua,
	}).Error
	if err != nil {
		logrus.Errorf("创建登录日志失败")
		return
	}
}
