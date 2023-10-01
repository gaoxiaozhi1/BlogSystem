package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
)

// 系统路由信息
func SettingsRouter(router *gin.RouterGroup) {
	settingsApi := api.ApiGroupApp.SettingsApi
	router.GET("settings/:name", settingsApi.SettingsInfoView)
	router.PUT("settings/:name", settingsApi.SettingsInfoUpdateView)

}
