package middleware

import (
	"blogX_server/common/resp"
	"blogX_server/utils/email_store"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
)

type EmailVerifyMiddlewareRequest struct {
	EmailID   string `json:"email_id" binding:"required"`
	EmailCode string `json:"email_code" binding:"required"`
}

// EmailVerifyMiddleware 邮箱验证码验证
func EmailVerifyMiddleware(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		resp.FailWithMsg("获取请求体错误", c)
		c.Abort()
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	var cr EmailVerifyMiddlewareRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		logrus.Errorf("邮箱验证失败 %s", err)
		resp.FailWithMsg("邮箱验证失败", c)
		c.Abort()
		return
	}
	info, ok := email_store.Verify(cr.EmailID, cr.EmailCode)
	if !ok {
		resp.FailWithMsg("邮箱验证码校验失败", c)
		c.Abort()
		return
	}
	c.Set("email", info.Email)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
}
