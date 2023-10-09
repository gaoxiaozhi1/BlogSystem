package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
)

func MenuRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.MenuApi
	router.POST("menus", app.MenuCreateView)
	router.GET("menus", app.MenuListView)
	router.GET("menu_names", app.MenuNameListView)
	router.GET("menus/:id", app.MenuDetailView)
	router.PUT("menus/:id", app.MenuUpdateView)
	router.DELETE("menus", app.MenuRemoveView)
}
