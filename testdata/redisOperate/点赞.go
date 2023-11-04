package main

import (
	"fmt"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/service/redis_ser"
)

func main() {
	core.InitConf()
	global.Log = core.InitLogger()
	global.Redis = core.ConnectRedis()
	redis_ser.Digg("dvO-g4sBc9tzXF2QLZIy")
	fmt.Println(redis_ser.GetDigg("XXX"))
	fmt.Println(redis_ser.GetDiggInfo())
	//redis_ser.DiggClear()
}
