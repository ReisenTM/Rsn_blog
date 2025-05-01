package conf

import "fmt"

type DB struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Debug    bool   `yaml:"debug"`  //是否打印全部日志
	Source   string `yaml:"source"` //数据库的源 pgsql mysql
}

// Dsn 拼接dsn
func (d *DB) Dsn() string {
	//拼接dsn
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		d.Host, d.Port, d.User, d.Password, d.Database)
	return dsn
}

// Empty 判空
func (d *DB) Empty() bool {
	return d.User == "" && d.Password == "" && d.Host == "" && d.Database == "" && d.Source == ""
}
