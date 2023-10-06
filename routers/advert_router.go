package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
)

func AdvertRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.AdvertAPI
	router.POST("adverts", app.AdvertCreate)
}
