package flags

import (
	"blogX_server/Model"
	"blogX_server/global"
	"github.com/sirupsen/logrus"
)

func FlagDB() {
	err := global.DB.AutoMigrate(
		&Model.UserModel{},
		&Model.UserConfModel{},
		&Model.ArticleModel{},
		&Model.CategoryModel{},
		&Model.ArticleLikeModel{},
		&Model.CollectModel{},
		&Model.UserArticleCollectModel{},
		&Model.UserArticleHistoryModel{},
		&Model.CommentModel{},
		&Model.BannerModel{},
		&Model.LogModel{},
		&Model.GlobalNotificationModel{},
	)
	if err != nil {
		logrus.Errorf("自动迁移失败 %s", err)
		return
	}
	logrus.Infof("自动迁移成功")
}
