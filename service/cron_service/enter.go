package cron_service

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"time"
)

func Cron() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))

	// 每天2点去同步文章数据
	_, err := crontab.AddFunc("0 0 2 * * *", SyncArticle)
	if err != nil {
		logrus.Errorf("crontab.AddFunc err:%v", err)
		return
	}

	crontab.Start()
}
