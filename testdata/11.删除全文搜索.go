package main

import (
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/service/es_ser"
)

func main() {
	// 读取配置文件
	core.InitConf()
	// 初始化日志
	global.Log = core.InitLogger()
	// 连接es
	global.ESClient = core.EsConnect()
	
	es_ser.DeleteFullTextByArticleID("e_MzlYsBc9tzXF2QOJKr")
}
