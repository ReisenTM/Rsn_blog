package flags

import (
	"blogX_server/global"
	"blogX_server/model"
	"github.com/sirupsen/logrus"
)

// FlagDB 迁移数据库
func FlagDB() {

	err := global.DB.AutoMigrate(
		&model.UserModel{},
		&model.UserConfModel{},
		&model.ArticleModel{},
		&model.CategoryModel{},
		&model.CollectModel{},
		&model.UserArticleCollectModel{},
		&model.UserArticleHistoryModel{},
		&model.CommentModel{},
		&model.BannerModel{},
		&model.LogModel{},
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
		&model.UserFocusModel{},              //好友关注列表
	)
	if err != nil {
		logrus.Errorf("自动迁移失败 %s", err)
		return
	}
	logrus.Infof("自动迁移成功")
}
