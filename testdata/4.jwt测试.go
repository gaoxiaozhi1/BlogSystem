package main

import (
	"fmt"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/utils/jwts"
)

func main() {
	// 读取配置文件
	core.InitConf()
	// 初始化日志
	global.Log = core.InitLogger()

	token, err := jwts.GenToken(jwts.JwtPayLoad{
		//Username: "admin",
		NickName: "admin",
		Role:     1,
		UserID:   1,
	})
	fmt.Println(token, err)

	// 两个小时过期一次
	claims, err := jwts.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwibmlja19uYW1lIjoiYWRtaW4iLCJyb2xlIjoxLCJ1c2VyX2lkIjoxLCJleHAiOjE2OTY4Mzk5NTkuNzA4ODI1fQ.ykRjKIrNOSv8_MN_Uvz25FXIaQ9DR0BNbBKVQGTG_kc")
	fmt.Println(claims, err)
}
