package user_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/models"
	"gvb_server/models/ctype"
	"gvb_server/models/res"
	"gvb_server/service/common"
	"gvb_server/utils/desens"
	"gvb_server/utils/jwts"
)

// UserListView 用户列表页，也要分页
func (UserApi) UserListView(c *gin.Context) {
	// 获得解析token得到的东西
	_claims, _ := c.Get("claims")
	// 断言 (注意_claims的类型)
	claims := _claims.(*jwts.CustomClaims)

	var page models.PageInfo
	if err := c.ShouldBindQuery(&page); err != nil {
		res.FailWitheCode(res.ArgumentError, c) // 参数错误
		return
	}
	// 分页参数
	// 要区分游客和管理员的界面
	// 密码不能展示，加json:"-"就可以
	// 手机号和邮箱需要脱敏
	// user_name要么用星号，要么就不显示，管理员可以看到这个
	list, count, _ := common.ComList(models.UserModel{}, common.Option{
		PageInfo: page,
		Debug:    true,
	})

	var users []models.UserModel
	for _, user := range list {
		if ctype.Role(claims.Role) != ctype.PermissionAdmin {
			// 非管理员
			user.UserName = ""
		}
		// 脱敏
		user.Tel = desens.DesensitizationTel(user.Tel)
		user.Email = desens.DesensitizationEmail(user.Email)
		users = append(users, user)
	}

	res.OKWithList(users, count, c)
}
