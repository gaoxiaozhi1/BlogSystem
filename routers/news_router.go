package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
)

func NewsRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.NewsApi
	router.POST("news", app.NewListView)
}
