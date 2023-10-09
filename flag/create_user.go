package flag

import (
	"fmt"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/ctype"
	"gvb_server/utils/pwd"
)

// CreateUser 创建用户
// permissins 权限
func CreateUser(permissions string) {
	// 创建用户的逻辑
	// 用户名 昵称 密码 确认密码 邮箱
	var (
		userName   string
		nickName   string
		password   string
		rePassword string
		email      string
	)
	fmt.Printf("请输入用户名")
	fmt.Scan(&userName)
	fmt.Printf("请输入昵称")
	fmt.Scan(&nickName)
	fmt.Printf("请输入密码")
	fmt.Scan(&password)
	fmt.Printf("请再次输入密码")
	fmt.Scan(&rePassword)
	fmt.Printf("请输入邮箱")
	fmt.Scan(&email) // 不用必须输入;因为Scanln问题，所以还是改成必填的Scan

	// 判断用户名是否存在
	var userModel models.UserModel
	err := global.DB.Find(&userModel, "user_name = ?", userName).Error
	if err == nil {
		// 存在
		global.Log.Error("用户名已存在，请重新输入")
		return
	}

	// 校验两次密码
	if password != rePassword {
		global.Log.Error("两次密码不一致，请重新输入")
		return
	}

	// 对密码进行加密（hash）
	hashPwd := pwd.HashPwd(password)
	// 还可以用正则判断密码的强度，还可以用洛必达判断是否属于弱密码

	// 头像问题
	// 1. 默认头像
	// 2. 随机选择头像
	advatar := "/upload/avatars/default.png"

	// 入库
	role := ctype.PermissionUser // 默认是普通用户
	if permissions == "admin" {
		role = ctype.PermissionAdmin
	}

	err = global.DB.Create(&models.UserModel{
		UserName:   userName,
		NickName:   nickName,
		Password:   hashPwd,
		Email:      email,
		Role:       role,
		Avatar:     advatar,
		IP:         "127.0.0.1",
		Addr:       "内网地址",
		SignStatus: ctype.SignEmail,
	}).Error
	if err != nil {
		global.Log.Error(err)
		return
	}
	global.Log.Infof("用户 %s 创建成功", userName)
}
