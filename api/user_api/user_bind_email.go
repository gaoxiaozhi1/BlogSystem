package user_api

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/plugins/email"
	"gvb_server/utils/jwts"
	"gvb_server/utils/pwd"
	"gvb_server/utils/random"
)

type BindEmailRequest struct {
	Email    string  `json:"email" binding:"required,email" msg:"邮箱非法"`
	Code     *string `json:"code"`     // 验证码(方便用户判断用户是否传验证码)
	Password string  `json:"password"` // 邮箱登录的密码是另外设置的
}

// UserBindEmailView 用户绑定邮箱
func (UserApi) UserBindEmailView(c *gin.Context) {
	// 获得解析token得到的东西
	_claims, _ := c.Get("claims")
	// 断言 (注意_claims的类型)
	claims := _claims.(*jwts.CustomClaims)
	// 第一次输入是 邮箱
	// 后台会给这个邮箱发验证码
	var cr BindEmailRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	// 初始化session
	session := sessions.Default(c)
	if cr.Code == nil {
		// 第一次，后台发验证码
		// 验证码要和本次会话保持一致 --- session
		// 生成4位验证码(公共方法)，将生成的验证码存入session，要不然用户下一次用的时候，就不知道code是发给谁啦
		code := random.Code(4)
		// 写入session，可以保持会话
		// 初始化，然后就可以写东西啦
		session.Set("valid_code", code)
		err = session.Save()
		if err != nil {
			global.Log.Error(err)
			res.FailWithMessage("session错误", c)
			return
		}
		err = email.NewCode().Send(cr.Email, "你的验证码是 "+code)
		// 如果这里想处理快一些的话，就加一个go，即
		// go email.NewCode().Send(cr.Email, "你的验证码是 "+code)
		if err != nil {
			global.Log.Error(err)
			fmt.Println(err)
		}
		res.OKWithMessage("验证码已发送，请查收", c)
		return
	}
	// 第二次，用户输入邮箱，验证码，密码
	code := session.Get("valid_code")
	// 校验验证码
	if code != *cr.Code {
		res.FailWithMessage("验证码错误", c)
		return
	}
	// 查找当前用户是否存在
	var user models.UserModel
	err = global.DB.Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMessage("用户不存在", c)
		return
	}
	if len(cr.Password) < 4 {
		res.FailWithMessage("密码强度太低", c)
		return
	}
	hashPwd := pwd.HashPwd(cr.Password)
	// 用户绑定邮箱(第一次的邮箱和第二次的邮箱也要做一致性校验)
	// 就是如果收验证码是一个邮箱，绑定验证码又是另一个邮箱，那就很危险，到时候要验证一下，这里后续完善
	err = global.DB.Model(&user).Updates(map[string]any{
		"email":    cr.Email,
		"password": hashPwd,
	}).Error
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("绑定邮箱失败", c)
		return
	}
	// 完成绑定
	res.OKWithMessage("邮箱绑定成功", c)
	return
}
