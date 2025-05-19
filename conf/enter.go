package conf

type Config struct {
	System System `yaml:"system"`
	Log    Log    `yaml:"log"`
	DB     DB     `yaml:"db"`    //读库
	DB1    DB     `yaml:"db1"`   //写库
	Jwt    Jwt    `yaml:"jwt"`   //JWT
	Redis  Redis  `yaml:"redis"` //redis
	Site   Site   `yaml:"site"`
	QiNiu  QiNiu  `yaml:"qiniu"`
	AI     AI     `yaml:"ai"`
	QQ     QQ     `yaml:"qq"`
	Email  Email  `yaml:"email"`
	Upload Upload `yaml:"upload"`
}
