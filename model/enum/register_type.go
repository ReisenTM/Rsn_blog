package enum

type RegisterSourceType int8

// 注册来源
const (
	RegisterEmailSourceType    RegisterSourceType = 1
	RegisterQQSourceType       RegisterSourceType = 2
	RegisterTerminalSourceType RegisterSourceType = 3 //命令行创建
)
