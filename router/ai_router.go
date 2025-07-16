package router

import (
	"blogX_server/api"
	"blogX_server/api/ai_api"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
)

func AIRouter(r *gin.RouterGroup) {
	app := api.App.AIApi
	r.POST("ai/analysis", middleware.AuthMiddleware, middleware.BindJsonMiddleware[ai_api.ArticleAnalysisRequest], app.ArticleAnalysisView)
	r.GET("ai/article", middleware.AuthMiddleware, middleware.BindQueryMiddleware[ai_api.ArticleAiRequest], app.ArticleAiView)
}
