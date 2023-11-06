package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
)

func ChatRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.ChatApi
	router.GET("chat_groups", app.ChatGroupView)
}
