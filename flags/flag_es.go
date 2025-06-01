package flags

import (
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/service/es_service"
	"github.com/sirupsen/logrus"
)

func EsIndex() {
	if global.EsClient == nil {
		logrus.Warnf("未开启es连接")
		return
	}
	article := model.ArticleModel{}
	es_service.CreateIndexV2(article.Index(), article.Mapping())

}
