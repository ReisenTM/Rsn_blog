package flags

import (
	"blogX_server/global"
	"blogX_server/model"
	"github.com/sirupsen/logrus"
)

// FlagDB 迁移数据库
func FlagDB() {

	err := global.DB.AutoMigrate(
		&model.UserModel{},                   //用户表
		&model.UserConfModel{},               //用户配置表
		&model.ArticleModel{},                //文章表
		&model.CategoryModel{},               //分类表
		&model.CollectModel{},                //收藏夹表
		&model.UserArticleCollectModel{},     //用户收藏关系表
		&model.UserArticleHistoryModel{},     //浏览记录表
		&model.CommentModel{},                //评论表
		&model.BannerModel{},                 //封面表
		&model.LogModel{},                    //日志表
		&model.GlobalNotificationModel{},     //全局通知表
		&model.ImageModel{},                  //图片表
		&model.UserLoginModel{},              //登录表
		&model.UserTopArticleModel{},         //置顶文章
		&model.UserArticleFavorModel{},       //点赞表
		&model.UserCommentFavorModel{},       //评论点赞表
		&model.MessageModel{},                //消息表
		&model.UserMessageConfModel{},        //消息设置表
		&model.GlobalNotificationModel{},     //用户全局消息表
		&model.UserGlobalNotificationModel{}, //用户全局消息行为表
		&model.UserFocusModel{},              //好友关注表
		&model.ChatModel{},                   //聊天信息表
		&model.UserChatActionModel{},         //聊天行为表
	)
	if err != nil {
		logrus.Errorf("自动迁移失败 %s", err)
		return
	}
	logrus.Infof("自动迁移成功")
}
