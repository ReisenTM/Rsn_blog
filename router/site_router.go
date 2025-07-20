package router

import (
	"blogX_server/api"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	app := api.App.SiteApi
	r.GET("site/:name", app.SiteInfoView)
	r.GET("site/ai_conf", app.SiteInfoAIView)

	r.GET("site/qq_url", app.SiteInfoQQView)
	//更新需要管理员权限
	r.PUT("site/:name", middleware.AdminMiddleware, app.SiteUpdateView)
	return
}
