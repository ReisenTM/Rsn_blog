package flags

import (
	"blogX_server/global"
	"blogX_server/model"
	"github.com/sirupsen/logrus"
)

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
		&model.GlobalNotificationModel{},
	)
	if err != nil {
		logrus.Errorf("自动迁移失败 %s", err)
		return
	}
	logrus.Infof("自动迁移成功")
}
