package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/service/redis_service/redis_jwt"
	"github.com/gin-gonic/gin"
)

// LogoutView 注销
func (UserApi) LogoutView(c *gin.Context) {
	token := c.Request.Header.Get("token")
	redis_jwt.RedisBlackList(token, redis_jwt.UserBanType)

	resp.OKWithMsg("注销成功", c)
}
