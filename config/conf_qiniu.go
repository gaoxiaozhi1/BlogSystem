package config

type QiNiu struct {
	Enable    bool    `json:"enable" yaml:"enable"`         // 是否启用七牛云存储
	AccessKey string  `json:"access_key" yaml:"access_key"` // 用于标识用户，用户将 AccessKey 放入访问请求，以便七牛云存储识别访问者的身份。
	SecretKey string  `json:"secret_key" yaml:"secret_key"` // 是用于加密签名字符串和服务器端验证签名字符串的密钥。
	Bucket    string  `json:"bucket" yaml:"bucket"`         // 存储桶的名字
	CDN       string  `json:"cdn" yaml:"cdn"`               // 访问图片的地址的前缀
	Zone      string  `json:"zone" yaml:"zone"`             // 存储的地区
	Size      float64 `json:"size" yaml:"size"`             // 存储的大小限制，单位是MB(超过5MB的就不存啦)
}
