package router

import (
	"blogX_server/api"
	"blogX_server/api/focus_api"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
)

func FocusRouter(r *gin.RouterGroup) {
	app := api.App.FocusApi
	r.POST("focus", middleware.AuthMiddleware, middleware.BindJsonMiddleware[focus_api.FocusUserRequest], app.FocusUserView)
	r.DELETE("focus", middleware.AuthMiddleware, middleware.BindJsonMiddleware[focus_api.FocusUserRequest], app.UnFocusUserView)
	r.GET("focus/my_focus", middleware.BindQueryMiddleware[focus_api.FocusUserListRequest], app.FocusUserListView)
	r.GET("focus/my_fans", middleware.BindQueryMiddleware[focus_api.FocusUserListRequest], app.FansUserListView)
}
