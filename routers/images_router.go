package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
)

func ImagesRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.ImagesApi
	router.POST("images", app.ImageUploadView)

}
