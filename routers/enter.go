package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
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

	router.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	// 测试请求接口，可以用于测试
	//router.GET("/login", user_api.UserApi{}.QQLoginView)
	// 在浏览器访问：blog.fengfengzhidao.com/login?flag=qq&code=057F987F9A27DFOD6C1E210552667A87，就会被返回code

	// 传入路由分组
	apiRouterGroup := router.Group("api")

	// 路由分层
	SettingsRouter(apiRouterGroup) // 系统配置api
	ImagesRouter(apiRouterGroup)   // 图片管理api
	AdvertRouter(apiRouterGroup)   // 广告管理api
	MenuRouter(apiRouterGroup)     // 菜单管理api
	UserRouter(apiRouterGroup)     // 用户管理api
	TagRouter(apiRouterGroup)      // 标签管理api
	MessageRouter(apiRouterGroup)  // 消息管理api
	ArticleRouter(apiRouterGroup)  // 广告管理api
	return router
}
