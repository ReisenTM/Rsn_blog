package middleware

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
)

type CaptchaMiddlewareRequest struct {
	CaptchaId   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
}

func CaptchaMiddleware(c *gin.Context) {
	if !global.Config.Site.Login.Captcha {
		//如果没开验证码
		return
	}
	body, err := c.GetRawData()
	if err != nil {
		resp.FailWithMsg("请求错误", c)
		c.Abort()
		return
	}
	//用掉body，写回去
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	var req CaptchaMiddlewareRequest
	err = c.ShouldBindJSON(&req)
	
	if err != nil {
		logrus.Errorf("图形验证失败 %s", err)
		resp.FailWithError(err, c)
		c.Abort()
		return
	}
	if !global.CaptchaStore.Verify(req.CaptchaId, req.CaptchaCode, true) {
		resp.FailWithMsg("验证码验证失败", c)
		c.Abort()
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
}
