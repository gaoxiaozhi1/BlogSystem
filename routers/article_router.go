package routers

import (
	"github.com/gin-gonic/gin"
	"gvb_server/api"
	"gvb_server/middleware"
)

func ArticleRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.ArticleApi
	router.POST("articles", middleware.JwtAuth(), app.ArticleCreateView)
	router.GET("articles", app.ArticleListView)
	router.GET("articles/detail", app.ArticleDetailByTitleView)
	router.GET("articles/calendar", app.ArticleCalendarView)
	router.GET("articles/tags", app.ArticleTagListView)
	router.GET("articles/:id", app.ArticleDetailView)
	router.PUT("articles", app.ArticleUpdateView)
	router.DELETE("articles", app.ArticleRemoveView)
	router.POST("articles/collects", middleware.JwtAuth(), app.ArticleCollCreateView)
	router.GET("articles/collects", middleware.JwtAuth(), app.ArticleCollListView)
	router.DELETE("articles/collects", middleware.JwtAuth(), app.ArticleCollBatchRemoveView)
	router.GET("articles/text", app.FullTextSearchView) // 全文搜索
}
