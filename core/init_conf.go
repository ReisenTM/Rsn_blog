package core

import (
	"blogX_server/conf"
	"blogX_server/flags"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// ReadConf 读取配置文件
func ReadConf() *conf.Config {
	byteData, err := os.ReadFile(flags.FlagOptions.File)
	if err != nil {
		panic(err)
	}
	c := new(conf.Config)
	err = yaml.Unmarshal(byteData, c)
	if err != nil {
		panic(fmt.Errorf("yaml配置文件格式错误%v", err))
	}
	fmt.Printf("读取配置文件成功:%s\n", flags.FlagOptions.File)
	return c
}
