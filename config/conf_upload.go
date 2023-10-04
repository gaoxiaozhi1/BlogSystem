package config

// 本地图片
type Upload struct {
	Size int    `yaml:"size" json:"size"` // 图片上传的大小，MB
	Path string `yaml:"path" json:"path"` // 图片上传的目录
}
