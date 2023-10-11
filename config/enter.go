package config

// 第一件事情就是执行这里面的事情
type Config struct {
	Mysql    Mysql    `yaml:"mysql"`
	Logger   Logger   `yaml:"logger"`
	System   System   `yaml:"system"`
	Upload   Upload   `yaml:"upload"`
	SiteInfo SiteInfo `json:"site_info"`
	QQ       QQ       `yaml:"qq"`
	QiNiu    QiNiu    `yaml:"qi_niu"`
	Email    Email    `yaml:"email"`
	Jwy      Jwy      `yaml:"jwy"`
	Redis    Redis    `yaml:"redis"`
}
