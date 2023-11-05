package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
	"gvb_server/middleware"
)

func CommentRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.CommentApi
	router.POST("comments", middleware.JwtAuth(), app.CommentCreateView)
	router.GET("comments", app.CommentListView)
	router.GET("comments/:id", app.CommentDiggView)
	router.DELETE("comments/:id", app.CommentRemoveView)
}
