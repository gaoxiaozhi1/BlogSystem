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

// Reverse 任意类型的切片反转
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}
