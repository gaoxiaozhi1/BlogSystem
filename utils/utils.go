package utils

// Inlist 判断某字符串是否存在在字符串列表里面
func Inlist(key string, list []string) bool {
	for _, s := range list {
		if s == key {
			return true
		}
	}
	return false
}
