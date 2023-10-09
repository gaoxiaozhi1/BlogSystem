package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
)

func UserRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.UserApi
	router.POST("email_login", app.EmailLoginView)
}
