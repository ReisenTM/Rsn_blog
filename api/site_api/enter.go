package site_api

import (
	"blogX_server/common/resp"
	"blogX_server/conf"
	"blogX_server/core"
	"blogX_server/global"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
)

// SiteApi 站点管理Api入口
type SiteApi struct {
}
type SiteInfoRequest struct {
	Name string `uri:"name"`
}

// SiteInfoView 站点视图管理
func (sa *SiteApi) SiteInfoView(c *gin.Context) {
	var req SiteInfoRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	if req.Name == "site" {
		global.Config.Site.About.Version = global.Version
		resp.OkWithData(global.Config.Site, c)
		return
	}

	// 判断角色是不是管理员
	middleware.AdminMiddleware(c)
	_, ok := c.Get("claims")
	if !ok {
		return
	}
	var respData any
	switch req.Name {
	//不能暴露给前端的配置，硬编码
	case "email":
		res := global.Config.Email
		res.AuthCode = "******"
		respData = res
	case "qq":
		res := global.Config.QQ
		res.AppKey = "******"
		respData = res
	case "qiNiu":
		res := global.Config.QiNiu
		res.SecretKey = "******"
		respData = res
	case "ai":
		res := global.Config.AI
		res.SecretKey = "******"
		respData = res
	default:
		resp.FailWithMsg("不存在的配置", c)
		return
	}
	resp.OkWithData(respData, c)
	return
}

func (sa *SiteApi) SiteInfoQQView(c *gin.Context) {
	resp.OkWithData(global.Config.QQ.Url(), c)
}

type SiteUpdateRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required"`
}

// SiteUpdateView 更新站点配置
func (sa *SiteApi) SiteUpdateView(c *gin.Context) {
	//log := log_service.GetLog(c)
	var req SiteInfoRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	//根据请求判断哪个配置要修改
	var res any
	switch req.Name {
	case "site":
		var siteConf conf.Site
		err = c.ShouldBindJSON(&siteConf)
		res = siteConf
	case "email":
		var emailConf conf.Email
		err = c.ShouldBindJSON(&emailConf)
		res = emailConf
	case "qq":
		var qqConf conf.QQ
		err = c.ShouldBindJSON(&qqConf)
		res = qqConf
	case "ai":
		var aiConf conf.AI
		err = c.ShouldBindJSON(&aiConf)
		res = aiConf
	case "qiNiu":
		var qiniuConf conf.QQ
		err = c.ShouldBindJSON(&qiniuConf)
		res = qiniuConf
	default:
		resp.FailWithMsg("不存在的配置", c)
		return
	}
	//统一处理错误
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	//根据断言判断哪个配置被修改了,并保存到内存中
	switch s := res.(type) {
	case conf.Site:
		// 判断站点信息前后关联性是否影响
		err = UpdateSite(s)
		if err != nil {
			resp.FailWithError(err, c)
			return
		}
		global.Config.Site = s
	case conf.Email:
		if s.AuthCode == "******" {
			s.AuthCode = global.Config.Email.AuthCode
		}
		global.Config.Email = s
	case conf.QQ:
		if s.AppKey == "******" {
			s.AppKey = global.Config.QQ.AppKey
		}
		global.Config.QQ = s
	case conf.QiNiu:
		if s.SecretKey == "******" {
			s.SecretKey = global.Config.QiNiu.SecretKey
		}
		global.Config.QiNiu = s
	case conf.AI:
		if s.SecretKey == "******" {
			s.SecretKey = global.Config.AI.SecretKey
		}
		global.Config.AI = s
	}
	//配置文件持久化
	core.SaveConf()
	resp.OKWithMsg("更新站点配置成功", c)
	return
}

// UpdateSite 判断配置之间关联关系
func UpdateSite(site conf.Site) error {
	//TODO
	return nil
}
