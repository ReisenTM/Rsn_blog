package core

import (
	"blogX_server/global"
	"blogX_server/service/river_service"
	"github.com/sirupsen/logrus"
)

func InitMysqlEs() {
	if global.Config.Es.Addr == "" {
		logrus.Infof("未配置es,关闭同步")
		return
	}
	r, err := river_service.NewRiver()
	if err != nil {
		logrus.Fatalf("初始化mysql-es err: %v", err)
	}
	go r.Run() //启动同步程序
}
