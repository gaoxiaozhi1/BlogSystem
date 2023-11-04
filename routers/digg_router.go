package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
)

func DiggRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.DiggApi
	router.POST("digg/article", app.DiggArticleView)
}
