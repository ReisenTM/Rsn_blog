package router

import (
	"blogX_server/api"
	"blogX_server/api/chat_api"
	"blogX_server/middleware"
	"blogX_server/model"
	"github.com/gin-gonic/gin"
)

func ChatRouter(r *gin.RouterGroup) {
	app := api.App.ChatApi
	r.GET("chat", middleware.AuthMiddleware, middleware.BindQueryMiddleware[chat_api.ChatListRequest], app.ChatListView)
	r.GET("chat/session", middleware.AuthMiddleware, middleware.BindQueryMiddleware[chat_api.SessionListRequest], app.SessionListView)
	r.DELETE("chat", middleware.AuthMiddleware, middleware.BindJsonMiddleware[model.RemoveRequest], app.UserChatDeleteView)
	r.DELETE("chat/user/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[model.IDRequest], app.UserChatDeleteByUserView)
	r.POST("chat/read/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[model.IDRequest], app.ChatReadView)
	r.GET("chat/ws", app.ChatView)
}
