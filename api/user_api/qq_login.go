package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/ctype"
	"gvb_server/models/res"
	"gvb_server/plugins/qq"
	"gvb_server/utils/jwts"
	"gvb_server/utils/pwd"
	"gvb_server/utils/random"
)

func (UserApi) QQLoginView(c *gin.Context) {
	// 先得到code
	code := c.Query("code")
	if code == "" {
		res.FailWithMessage("没有code", c)
		return
	}
	fmt.Println(code)
	// QQ登录 --- 拿到openID
	qqInfo, err := qq.NewQQLogin(code) // 这里只能调一次，只能请求一次，不然code是失效的也没用
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}
	openID := qqInfo.OpenID
	// 根据openID判断用户是否存在
	var user models.UserModel
	// token 其他平台的唯一id
	err = global.DB.Take(&user, "token = ?", openID).Error
	if err != nil {
		// 不存在，就注册
		hashPwd := pwd.HashPwd(random.RandString(16))
		user = models.UserModel{
			NickName:   qqInfo.Nickname,
			UserName:   openID,        // qq登录，邮箱+密码，所以不用弄UserName登录，直接openid就行
			Password:   hashPwd,       // 随机生成16位密码
			Avatar:     qqInfo.Avatar, // 头像
			Addr:       "内网",          // 根据ip算地址，目前没有公网ip，所以之后再讲
			Token:      openID,
			IP:         c.ClientIP(),
			Role:       ctype.PermissionUser,
			SignStatus: ctype.SignQQ,
		}
		err = global.DB.Create(&user).Error
		if err != nil {
			global.Log.Error(err)
			res.FailWithMessage("注册失败", c)
			return
		}
		// 注册之后要把token返回给前端，前端才能做登录操作，才能做token的持久化，才能登录
	}

	// 登录操作
	// 登陆成功，生成token
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		NickName: user.NickName,
		Role:     int(user.Role),
		UserID:   user.ID,
	})

	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("token生成失败", c)
		return
	}

	res.OKWithData(token, c)
}
