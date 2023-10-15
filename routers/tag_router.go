package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
)

func TagRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.TagApi
	router.POST("tags", app.TagCreateView)
	router.GET("tags", app.TagListView)
	router.PUT("tags/:id", app.TagUpdateView)
	router.DELETE("tags", app.TagRemoveView)
}
