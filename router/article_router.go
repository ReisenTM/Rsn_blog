package router

import (
	"blogX_server/api"
	"blogX_server/api/article_api"
	"blogX_server/middleware"
	"blogX_server/model"
	"github.com/gin-gonic/gin"
)

func ArticleRouter(r *gin.RouterGroup) {
	app := api.App.ArticleApi
	//文章
	r.POST("article", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleCreateRequest], app.CreateArticleView)
	r.GET("article", middleware.BindQueryMiddleware[article_api.ArticleListRequest], app.ArticleListView)
	r.PUT("article", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleUpdateRequest], app.UpdateArticleView)
	r.GET("article/:id", middleware.BindUriMiddleware[model.IDRequest], app.ArticleDetailView)
	r.DELETE("article/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[model.IDRequest], app.ArticleRemoveUserView)
	r.DELETE("article", middleware.AdminMiddleware, middleware.BindJsonMiddleware[model.RemoveRequest], app.ArticleRemoveAdminView)
	//个性化推荐
	r.GET("article/auth_recommend", middleware.BindQueryMiddleware[article_api.AuthRecommendRequest], app.AuthRecommendView)
	//审核
	r.POST("article/examine", middleware.AdminMiddleware, middleware.BindJsonMiddleware[article_api.ArticleExamineRequest], app.ArticleExamineView)
	//点赞
	r.GET("article/favor/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[model.IDRequest], app.ArticleFavorView)
	//收藏夹-文章
	r.POST("article/collect", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleCollectRequest], app.ArticleCollectView)
	r.DELETE("article/collect", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleCollectPatchRemoveRequest], app.ArticleCollectPatchRemoveView)

	//浏览记录
	r.POST("article/history", middleware.BindJsonMiddleware[article_api.ArticleLookRequest], app.ArticleLookView)
	r.GET("article/history", middleware.AuthMiddleware, middleware.BindQueryMiddleware[article_api.ArticleLookListRequest], app.ArticleLookListView)
	r.DELETE("article/history", middleware.AuthMiddleware, middleware.BindJsonMiddleware[model.RemoveRequest], app.ArticleLookRemoveView)
	//分类
	r.POST("category", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.ArticleCategoryRequest], app.CategoryCreateView)
	r.GET("category", middleware.BindQueryMiddleware[article_api.ArticleCategoryListRequest], app.CategoryListView)
	r.DELETE("category", middleware.AuthMiddleware, middleware.BindJsonMiddleware[model.RemoveRequest], app.CategoryRemoveView)
	//收藏夹
	r.POST("collect", middleware.AuthMiddleware, middleware.BindJsonMiddleware[article_api.CollectCreateRequest], app.CollectCreateView)
	r.GET("collect", middleware.BindQueryMiddleware[article_api.CollectListRequest], app.CollectListView)
	r.DELETE("collect", middleware.AuthMiddleware, middleware.BindJsonMiddleware[model.RemoveRequest], app.CollectRemoveView)
	//方便前端调用的Option接口
	r.GET("category/options", middleware.AuthMiddleware, app.CategoryOptionsView)
	r.GET("article/tag/options", middleware.AuthMiddleware, app.ArticleTagOptionsView)
}
