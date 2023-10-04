package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/global"
)

// 初始化路由
func InitRouter() *gin.Engine {
	/*
		简洁日志部分
		[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
		 - using env:   export GIN_MODE=release （在.yaml文件中）
		 - using code:  gin.SetMode(gin.ReleaseMode) （在路由初始化处）
	*/
	gin.SetMode(global.Config.System.Env)

	router := gin.Default()

	// 传入路由分组
	apiRouterGroup := router.Group("api")

	// 路由分层
	SettingsRouter(apiRouterGroup) // 系统配置api
	ImagesRouter(apiRouterGroup)   // 图片管理api
	return router
}
