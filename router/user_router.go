package router

import (
	"blogX_server/api"
	"blogX_server/api/user_api"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.GET("user", middleware.AdminMiddleware, middleware.BindQueryMiddleware[user_api.UserListRequest], app.UserListView)
	r.POST("user/send_email", middleware.CaptchaMiddleware, app.SendEmailView)
	r.POST("user/email", middleware.EmailVerifyMiddleware, app.RegisterEmailView)
	r.POST("user/login", middleware.CaptchaMiddleware, app.PwdLoginView)
	r.GET("user/detail", middleware.AuthMiddleware, app.UserDetailView)
	r.GET("user/login", middleware.AuthMiddleware, app.UserLoginListView)
	r.GET("user/info", app.UserBaseInfoView)
	r.PUT("user/password", middleware.AuthMiddleware, app.UpdatePasswordView)
	r.PUT("user/password/reset", middleware.EmailVerifyMiddleware, app.ResetPasswordView)
	r.PUT("user/email/bind", middleware.EmailVerifyMiddleware, middleware.AuthMiddleware, app.BindEmailView)
	r.PUT("user", middleware.AuthMiddleware, app.UserInfoUpdateView)
	r.PUT("user/admin", middleware.AdminMiddleware, app.AdminUserInfoUpdateView)
	r.POST("user/article/top", middleware.AuthMiddleware, middleware.BindJsonMiddleware[user_api.UserArticleTopRequest], app.UserArticleTopView)
	r.DELETE("user/logout", middleware.AuthMiddleware, app.LogoutView)
}
