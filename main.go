package main

import (
	"blogX_server/core"
	"blogX_server/flags"
	"blogX_server/global"
	"blogX_server/router"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitDefaultLogus()
	global.DB = core.InitDb()
	global.Redis = core.InitRedis()

	flags.Run()
	//启动程序
	router.Run()
}
