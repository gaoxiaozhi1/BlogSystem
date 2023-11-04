package redis_ser

import (
	"gvb_server/global"
	"gvb_server/utils"
	"time"
)

// 注意这里的是普通的函数，没有type RedisService struct{}这个东西
// 这里的函数会被user_ser和jwt_auth调用嗷
// 前缀
const prefix = "logout_"

// Logout 针对注销的操作
func Logout(token string, diff time.Duration) error {
	err := global.Redis.Set(prefix+token, "", diff).Err()
	return err
}

// CheckLogout 检测是否注销就是判断是否在redis中
func CheckLogout(token string) bool {
	keys := global.Redis.Keys(prefix + "*").Val()
	if utils.Inlist(prefix+token, keys) {
		return true // 在里面
	}
	return false // 不在里面
}
