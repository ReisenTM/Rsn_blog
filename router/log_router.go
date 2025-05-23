package router

import (
	"blogX_server/api"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
)

func LogRouter(rr *gin.RouterGroup) {
	app := api.App.LogApi
	r := rr.Group("").Use(middleware.AdminMiddleware)
	//日志都需要管理员权限
	r.GET("logs", app.LogListView)
	r.GET("logs/:id", app.LogReadView)
	r.DELETE("logs", app.LogDeleteView)
	return
}
