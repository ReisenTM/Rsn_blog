package router

import (
	"blogX_server/api"
	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	app := api.App.SiteApi
	r.GET("/site", app.SiteInfoView)
	r.PUT("/site", app.SiteUpdateView)
	return
}
