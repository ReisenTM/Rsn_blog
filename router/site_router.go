package router

import (
	"blogX_server/api"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	app := api.App.SiteApi
	r.GET("/site", app.SiteInfoView)
	//更新需要管理员权限
	r.PUT("/site", middleware.AdminMiddleware, app.SiteUpdateView)
	return
}
