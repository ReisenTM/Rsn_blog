package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/service/email_service"
	"blogX_server/utils/email_store"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
)

const (
	EmailRegType = iota + 1
	EmailResetType
)

type SendEmailRequest struct {
	Email string `json:"email" binding:"required"`
	Type  uint8  `json:"type" binding:"oneof=1 2"` //1注册 2 重置密码
}

type SendEmailResponse struct {
	CodeID string `json:"code_id"` //验证码id
}

// SendEmailView 发送操作提示邮件
func (UserApi) SendEmailView(c *gin.Context) {
	var cr SendEmailRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		resp.FailWithError(err, c)
		return
	}
	//生成验证码和id.存到验证码存储器中
	code := base64Captcha.RandText(4, "1234567890")
	id := base64Captcha.RandomId()
	var err error
	switch cr.Type {
	case EmailRegType:
		//先查邮箱是否存在
		var um model.UserModel
		err = global.DB.Take(&um, "email = ?", cr.Email).Error
		if err == nil && um.Email != "" {
			resp.FailWithMsg("邮箱已存在，请登录", c)
			return
		}
		err = email_service.SendRegCode(cr.Email, code)
	case EmailResetType:
		var user model.UserModel
		err = global.DB.Take(&user, "email = ?", cr.Email).Error
		if err != nil {
			resp.FailWithMsg("该邮箱不存在", c)
			return
		}
		// 还必须得是邮箱注册的

		err = email_service.SendResetCode(cr.Email, code)
	}
	if err != nil {
		logrus.Errorf("邮件发送失败 %s", err)
		resp.FailWithMsg("发送邮件失败,%s", c)
		return
	}
	//err = global.CaptchaStore.Set(id, code)
	email_store.Set(id, cr.Email, code)

	resp.OkWithData(SendEmailResponse{CodeID: id}, c)
}
