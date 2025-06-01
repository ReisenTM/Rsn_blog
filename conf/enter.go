package conf

type Config struct {
	System System `yaml:"system"`
	Log    Log    `yaml:"log"`
	DB     []DB   `yaml:"db"`     //连接的数据库
	Jwt    Jwt    `yaml:"jwt"`    //JWT
	Redis  Redis  `yaml:"redis"`  //redis
	Site   Site   `yaml:"site"`   //站点配置
	QiNiu  QiNiu  `yaml:"qiniu"`  //七牛云配置
	AI     AI     `yaml:"ai"`     //AI配置
	QQ     QQ     `yaml:"qq"`     //QQ配置
	Email  Email  `yaml:"email"`  //邮箱配置
	Upload Upload `yaml:"upload"` //图片后缀上传白名单
	Es     ES     `yaml:"es"`     //elastic配置
	River  River  `yaml:"river"`  //es-mysql同步
}
