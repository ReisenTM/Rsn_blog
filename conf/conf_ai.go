package conf

type AI struct {
	Enable    bool   `yaml:"enable" json:"enable"`
	SecretKey string `yaml:"secretKey" json:"secretKey"`
	Nickname  string `yaml:"nickname" json:"nickname"`
	Avatar    string `yaml:"avatar" json:"avatar"`
	Abstract  string `yaml:"abstract" json:"abstract"` //介绍
}
