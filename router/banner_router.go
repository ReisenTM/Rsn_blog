package router

import (
	"blogX_server/api"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
)

func BannerRouter(r *gin.RouterGroup) {
	app := api.App.BannerApi
	r.GET("banner", app.BannerListCreateView)
	r.POST("banner", middleware.AdminMiddleware, app.BannerCreateView)
	r.DELETE("banner", middleware.AdminMiddleware, app.BannerRemoveView)
	r.PUT("banner/:id", middleware.AdminMiddleware, app.BannerUpdateView)
}
