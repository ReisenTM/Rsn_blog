package router

import (
	"blogX_server/api"
	"blogX_server/api/site_msg_api"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
)

func SiteMsgRouter(r *gin.RouterGroup) {
	app := api.App.SiteMsgApi
	r.GET("site_msg", middleware.AuthMiddleware, middleware.BindQueryMiddleware[site_msg_api.SiteMsgListRequest], app.SiteMsgListView)
	//消息配置查看
	r.GET("site_msg/conf", middleware.AuthMiddleware, app.UserSiteMessageConfView)
	//消息配置更新
	r.PUT("site_msg/conf", middleware.AuthMiddleware, middleware.BindJsonMiddleware[site_msg_api.UserMessageConfUpdateRequest], app.UserSiteMessageConfUpdateView)
	//消息已读
	r.POST("site_msg", middleware.AuthMiddleware, middleware.BindJsonMiddleware[site_msg_api.SiteMsgReadRequest], app.SiteMsgReadView)
	//消息删除
	r.DELETE("site_msg", middleware.AuthMiddleware, middleware.BindJsonMiddleware[site_msg_api.SiteMsgRemoveRequest], app.SiteMsgRemoveView)
	//用户未读消息
	r.GET("site_msg/user", middleware.AuthMiddleware, app.UserMsgView)
}
