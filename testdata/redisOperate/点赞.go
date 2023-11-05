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
	redis_ser.NewDigg().Set("dvO-g4sBc9tzXF2QLZIy")
	fmt.Println(redis_ser.NewDigg().Get("dvO-g4sBc9tzXF2QLZIy"))
	fmt.Println(redis_ser.NewDigg().GetInfo())
	//redis_ser.DiggClear()
}
