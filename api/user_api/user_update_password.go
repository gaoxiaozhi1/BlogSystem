package user_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/utils/jwts"
	"gvb_server/utils/pwd"
)

type UpdatePasswordRequest struct {
	OldPwd string `json:"old_pwd" binding:"required" msg:"请输入旧密码"` // 旧密码
	Pwd    string `json:"pwd" binding:"required" msg:"请输入新密码"`     // 新密码
}

// UserUpdatePasswordView 修改登录人id
func (UserApi) UserUpdatePasswordView(c *gin.Context) {
	// 获取token解析后的claims
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	// 参数绑定
	var cr UpdatePasswordRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, &cr, c)
		return
	}

	// 查找用户是否存在
	var user models.UserModel
	err := global.DB.Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMessage("用户不存在", c)
		return
	}

	// 判断密码是否正确
	if !pwd.CheckPwd(user.Password, cr.OldPwd) {
		res.FailWithMessage("密码错误", c)
		return
	}
	// 旧密码正确就修改密码
	hashPwd := pwd.HashPwd(cr.Pwd)
	err = global.DB.Model(&user).Update("password", hashPwd).Error
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("密码修改失败", c)
		return
	}
	res.OKWithMessage("密码修改成功", c)
	return
}
