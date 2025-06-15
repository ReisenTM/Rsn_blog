package router

import (
	"blogX_server/api"
	"blogX_server/api/comment_api"
	"blogX_server/middleware"
	"blogX_server/model"
	"github.com/gin-gonic/gin"
)

func CommentRouter(r *gin.RouterGroup) {
	app := api.App.CommentApi
	r.POST("comment", middleware.AuthMiddleware, middleware.BindJsonMiddleware[comment_api.CommentCreateRequest], app.CommentCreateView)
	r.GET("comment/tree/:id", middleware.BindUriMiddleware[model.IDRequest], app.CommentTreeView)
	r.GET("comment", middleware.AuthMiddleware, middleware.BindQueryMiddleware[comment_api.CommentListRequest], app.CommentListView)
	r.DELETE("comment/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[model.IDRequest], app.CommentRemoveView)
	r.GET("comment/favor/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[model.IDRequest], app.CommentFavorView)
}
