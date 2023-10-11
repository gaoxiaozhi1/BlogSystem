package user_ser

import (
	"gvb_server/service/redis_ser"
	"gvb_server/utils/jwts"
	"time"
)

type UserService struct {
}

// Logout 注销操作（逻辑）:计算用户的token的过期时间，并将过期时间放入redis
func (UserService) Logout(claims *jwts.CustomClaims, token string) error {
	// 需要计算距离现在的过期时间(redis需要的过期时间是Duration类型)
	exp := claims.ExpiresAt // 最后的时间
	now := time.Now()       // 现在时间
	// 计算距离现在的过期时间
	diff := exp.Time.Sub(now) // 剩多少秒，单位是Duration

	// 操作redis
	return redis_ser.Logout(token, diff)
}
