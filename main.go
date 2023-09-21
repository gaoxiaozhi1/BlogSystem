package main

import (
	"gvb_server/core"
	"gvb_server/flag"
	"gvb_server/global"
	"gvb_server/routers"
)

func main() {
	// 读取配置文件
	core.InitConf()
	// 初始化日志
	global.Log = core.InitLogger()
	// 连接数据库
	global.DB = core.InitGorm()

	// 命令行参数绑定
	// 先迁移表结构
	option := flag.Parse()
	// 判断是否停止web项目
	if flag.IsWebStop(option) {
		flag.SwitchOption(option)
		return
	}

	// 初始化路由
	router := routers.InitRouter()
	addr := global.Config.System.Addr()
	global.Log.Infof("gvb_server运行在：%s", addr)

	// func (engine *Engine) Run(addr ...string) (err error) 所以要弄地址
	err := router.Run(addr)
	if err != nil {
		global.Log.Fatalf(err.Error())
	}

}
