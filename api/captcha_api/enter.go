package captcha_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
)

type CaptchaApi struct{}

type CaptchaResponse struct {
	CaptchaId string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

// CaptchaGenerateView 验证码生成接口
func (CaptchaApi) CaptchaGenerateView(c *gin.Context) {
	id, code, _, err := createCode()
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	resp.OkWithData(CaptchaResponse{
		CaptchaId: id,
		Captcha:   code,
	}, c)
}

// 根据自己需求更改验证码存储上限和过期时间
// var result = base64Captcha.NewMemoryStore(10240, 3*time.Minute)

// digitConfig 生成图形化数字验证码配置
func digitConfig() *base64Captcha.DriverDigit {
	digitType := &base64Captcha.DriverDigit{
		Height:   50,
		Width:    100,
		Length:   5,
		MaxSkew:  0.45,
		DotCount: 80,
	}
	return digitType
}

// CreateCode
// @Result id 验证码id
// @Result bse64s 图片base64编码
// @Result err 错误
func createCode() (string, string, string, error) {
	var driver base64Captcha.Driver
	//纯数字验证码
	driver = digitConfig()
	if driver == nil {
		logrus.Errorf("图形化数字验证码配置失败")
	}
	// 创建验证码并传入创建的类型的配置，以及存储的对象
	c := base64Captcha.NewCaptcha(driver, global.CaptchaStore)
	id, b64s, answer, err := c.Generate()
	return id, b64s, answer, err
}
