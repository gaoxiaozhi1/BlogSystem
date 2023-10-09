package jwts

import (
	"github.com/dgrijalva/jwt-go/v4"
	"gvb_server/global"
	"time"
)

// GenToken 创建Token
func GenToken(user JwtPayLoad) (string, error) {
	claim := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * time.Duration(global.Config.Jwy.Expires))), // 默认两个小时更新
			Issuer:    global.Config.Jwy.Issuer,                                                     // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	MySecret = []byte(global.Config.Jwy.Secret)
	return token.SignedString(MySecret)
}
