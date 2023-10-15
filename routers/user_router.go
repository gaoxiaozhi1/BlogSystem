package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gvb_server/api"
	"gvb_server/middleware"
)

var store = cookie.NewStore([]byte("secret")) // // 这里可以自己在配置文件 settings.yaml 中设置?

func UserRouter(router *gin.RouterGroup) {
	app := api.ApiGroupApp.UserApi
	// 这里session其实是一个中间件，那么我们同样可以自己写，自己写的话就可以控制session的过期时间
	router.Use(sessions.Sessions("sessionid", store))
	router.POST("email_login", app.EmailLoginView)
	router.POST("users", middleware.JwtAdmin(), app.UserCreateView)
	router.POST("login", app.QQLoginView) // 这里就是实际用的接口
	router.GET("users", middleware.JwtAuth(), app.UserListView)
	router.PUT("user_role", middleware.JwtAdmin(), app.UserUpdateRoleView)
	router.PUT("user_password", middleware.JwtAuth(), app.UserUpdatePasswordView)
	router.POST("logout", middleware.JwtAuth(), app.LogoutView)
	router.DELETE("users", middleware.JwtAdmin(), app.UserRemoveView)
	router.POST("user_bind_email", middleware.JwtAuth(), app.UserBindEmailView)
}
