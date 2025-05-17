package conf

// QiNiu 七牛云配置
type QiNiu struct {
	Enable    bool   `yaml:"enable" json:"enable"`
	AccessKey string `yaml:"accessKey" json:"accessKey"`
	SecretKey string `yaml:"secretKey" json:"secretKey"`
	Bucket    string `yaml:"bucket" json:"bucket"` //存储桶
	Uri       string `yaml:"uri" json:"uri"`
	Region    string `yaml:"region" json:"region"`
	Prefix    string `yaml:"prefix" json:"prefix"` //前缀
	Size      int    `yaml:"size" json:"size"`     // 大小限制 单位mb
	Expiry    int    `yaml:"expiry" json:"expiry"` // 过期时间 单位秒
}
