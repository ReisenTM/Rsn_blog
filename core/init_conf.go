package core

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var confPath = "settings.yaml"

type System struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type Config struct {
	System System `yaml:"system"`
}

// ReadConf 读取配置文件
func ReadConf() {
	byteData, err := os.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(byteData, &config)
	if err != nil {
		panic(fmt.Errorf("yaml配置文件格式错误%v", err))
	}
	fmt.Println(config)
}
