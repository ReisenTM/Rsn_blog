package flags

import (
	"flag"
	"os"
)

type Options struct {
	File    string
	DB      bool
	Version bool
}

var FlagOptions = new(Options)

// Parse flag绑定
func Parse() {
	flag.StringVar(&FlagOptions.File, "f", "settings.yaml", "配置文件")
	flag.BoolVar(&FlagOptions.DB, "db", false, "数据库迁移")
	flag.BoolVar(&FlagOptions.Version, "v", false, "版本")
	flag.Parse()
}

// Run flag实现
func Run() {
	if FlagOptions.DB {
		//执行数据库迁移
		FlagDB()
		os.Exit(0)
	}
}
