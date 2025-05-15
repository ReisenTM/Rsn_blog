package site_api

import (
	"blogX_server/common/resp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SiteApi 站点管理Api
type SiteApi struct {
}

// SiteInfoView 视图
func (siteApi *SiteApi) SiteInfoView(c *gin.Context) {
	resp.OkWithData("xxx", c)
	return
}

type SiteUpdateRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required"`
}

func (siteApi *SiteApi) SiteUpdateView(c *gin.Context) {
	//log := log_service.GetLog(c)
	var req SiteUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logrus.Errorf("参数绑定失败,%v", err)
		resp.FailWithError(err, c)
		return
	}
	resp.OkWithData(req, c)
	return
}
