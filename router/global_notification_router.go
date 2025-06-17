package router

import (
	"blogX_server/api"
	"blogX_server/api/global_notification_api"
	"blogX_server/middleware"
	"blogX_server/model"
	"github.com/gin-gonic/gin"
)

func GlobalNotificationRouter(r *gin.RouterGroup) {
	app := api.App.GlobalNotificationApi
	r.POST("global_notification", middleware.AdminMiddleware, middleware.BindJsonMiddleware[global_notification_api.CreateRequest], app.CreateView)
	r.GET("global_notification", middleware.AuthMiddleware, middleware.BindQueryMiddleware[global_notification_api.ListRequest], app.ListView)
	r.DELETE("global_notification", middleware.AdminMiddleware, middleware.BindJsonMiddleware[model.RemoveRequest], app.RemoveAdminView)
	r.POST("global_notification/user", middleware.AuthMiddleware, middleware.BindJsonMiddleware[global_notification_api.UserMsgActionRequest], app.UserMsgActionView)
}
