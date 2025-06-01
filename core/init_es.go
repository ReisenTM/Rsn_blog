package core

import (
	"blogX_server/global"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func InitES() *elastic.Client {
	es := global.Config.Es
	if !es.Enable || es.Addr == "" {
		logrus.Infof("未启用es连接")
		return nil
	}
	client, err := elastic.NewClient(
		elastic.SetURL(es.Addr),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(es.Username, es.Password),
	)
	if err != nil {
		logrus.Panicf("es连接失败 %s", err)
		return nil
	}
	logrus.Infof("es连接成功")
	return client
}
