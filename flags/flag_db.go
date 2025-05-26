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
		&model.ArticleLikeModel{},
		&model.CollectModel{},
		&model.UserArticleCollectModel{},
		&model.UserArticleHistoryModel{},
		&model.CommentModel{},
		&model.BannerModel{},
		&model.LogModel{},
		&model.GlobalNotificationModel{}, //全局通知表
		&model.ImageModel{},              //图片表
		&model.UserLoginModel{},          //登录表
	)
	if err != nil {
		logrus.Errorf("自动迁移失败 %s", err)
		return
	}
	logrus.Infof("自动迁移成功")
}
