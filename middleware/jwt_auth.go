package middleware

import (
	"github.com/gin-gonic/gin"
	"gvb_server/models/ctype"
	"gvb_server/models/res"
	"gvb_server/utils/jwts"
)

// 用户登录的中间件
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如何判断是管理员
		token := c.Request.Header.Get("token")
		if token == "" {
			res.FailWithMessage("未携带token", c)
			c.Abort() // 拦截
			return
		}
		// 解析token
		claims, err := jwts.ParseToken(token)
		if err != nil {
			res.FailWithMessage("token错误", c)
			c.Abort() // 拦截
			return
		}
		// 登录的用户
		c.Set("claims", claims) // 给需要使用中间件的函数取claims
	}
}

// JwtAdmin 管理员才有权限的中间件
func JwtAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如何判断是管理员
		token := c.Request.Header.Get("token")
		if token == "" {
			res.FailWithMessage("未携带token", c)
			c.Abort() // 拦截
			return
		}
		claims, err := jwts.ParseToken(token)
		if err != nil {
			res.FailWithMessage("token错误", c)
			c.Abort() // 拦截
			return
		}
		// 登录的用户
		if claims.Role != int(ctype.PermissionAdmin) {
			res.FailWithMessage("权限不足", c)
			c.Abort() // 拦截
			return
		}
		c.Set("claims", claims) // 给需要使用中间件的函数取claims
	}
}
