package router

import (
	"blogX_server/api"
	"blogX_server/api/search_api"
	"blogX_server/common"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
)

func SearchRouter(r *gin.RouterGroup) {
	app := api.App.SearchApi
	r.GET("search/article", middleware.BindQueryMiddleware[search_api.ArticleSearchRequest], app.ArticleSearchView)
	r.GET("search/tags", middleware.BindQueryMiddleware[common.PageInfo], app.TagAggView)
	r.GET("search/text", middleware.BindQueryMiddleware[search_api.TextSearchRequest], app.TextSearchView)
}
