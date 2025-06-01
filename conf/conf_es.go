package conf

import "fmt"

type ES struct {
	Addr     string `yaml:"url" json:"url"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	IsHttps  bool   `yaml:"is_https" json:"is_https"`
	Enable   bool   `yaml:"enable" json:"enable"`
}

func (e ES) Url() string {
	if e.IsHttps {
		return fmt.Sprintf("https://%s", e.Addr)
	}
	return fmt.Sprintf("http://%s", e.Addr)
}
