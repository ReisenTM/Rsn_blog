package global

import (
	"blogX_server/conf"
	"github.com/mojocn/base64Captcha"
	es "github.com/olivere/elastic/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Version 后端版本号
const Version = "1.0.1"

// global用来为项目内的方法提供操作对象
var (
	Config       *conf.Config                    //全局配置
	DB           *gorm.DB                        //数据库
	Redis        *redis.Client                   //redis
	CaptchaStore = base64Captcha.DefaultMemStore //图片验证码存储器
	EsClient     *es.Client                      //es
)
