package core

import (
	"blogX_server/global"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"time"
)

func InitDb() *gorm.DB {
	dc := global.Config.DB
	dc1 := global.Config.DB1

	if dc.Host == "" {
		logrus.Warnln("未配置数据库连接地址")
	}
	dsn := dc.Dsn()
	var myLogger logger.Interface
	//gorm连接
	if dc.Debug == true {
		//Debug环境下显示log
		myLogger = logger.Default.LogMode(logger.Info)
	} else {
		//正常模式下仅显示错误
		myLogger = logger.Default.LogMode(logger.Error)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, //不生成外键约束
		Logger:                                   myLogger,
	})
	if err != nil {
		logrus.Fatalf("连接数据库失败 %s", err)
	}
	//拿到原始sqlDB
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)  //最大空闲连接数
	sqlDB.SetMaxOpenConns(100) //最多可容纳
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infoln("数据库 连接成功")

	if !dc1.Empty() {
		//如果配置了读写库，就进行读写分离的注册
		err := db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{postgres.Open(dc1.Dsn())}, //写
			Replicas: []gorm.Dialector{postgres.Open(dc.Dsn())},  //读
			Policy:   dbresolver.RandomPolicy{},
		}))
		if err != nil {
			logrus.Fatal("读写分离配置错误\n", err)
		}
	}
	return db
}
