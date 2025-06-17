package router

import (
	"blogX_server/global"
	"blogX_server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run() {
	//设置运行模式
	gin.SetMode(global.Config.System.GinMode)

	r := gin.Default()
	//路径请求映射
	r.Static("/uploads", "uploads")
	//创建路由组
	nr := r.Group("/api")
	//使用中间件
	nr.Use(middleware.LogMiddleWare)

	SiteRouter(nr)
	LogRouter(nr)
	ImageRouter(nr)
	BannerRouter(nr)
	CaptchaRouter(nr)
	UserRouter(nr)
	ArticleRouter(nr)
	CommentRouter(nr)
	SiteMsgRouter(nr)
	//启动路由监听
	addr := global.Config.System.Addr()
	err := r.Run(addr)
	if err != nil {
		logrus.Errorf("server启动失败:%v", err)
		return
	}
}
