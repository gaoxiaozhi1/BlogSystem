package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
)

func AdvertRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.AdvertAPI
	router.POST("adverts", app.AdvertCreateView)
	router.GET("adverts", app.AdvertListView)
	router.PUT("adverts/:id", app.AdvertUpdateView)
	router.DELETE("adverts", app.AdvertRemoveView)
}
