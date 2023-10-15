package user_ser

import (
	"errors"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/ctype"
	"gvb_server/utils/pwd"
)

// 头像问题
// 1. 默认头像
// 2. 随机选择头像
const Avatar = "/upload/avatars/default.png"

func (UserService) CreateUser(userName, nickName, password string, role ctype.Role, email string, ip string) error {
	// 判断用户名是否存在
	var userModel models.UserModel
	err := global.DB.Take(&userModel, "user_name = ?", userName).Error
	if err == nil {
		return errors.New("用户名已存在")
	}

	// 对密码进行加密（hash）
	hashPwd := pwd.HashPwd(password)
	// 还可以用正则判断密码的强度，还可以用洛必达判断是否属于弱密码

	err = global.DB.Create(&models.UserModel{
		UserName:   userName,
		NickName:   nickName,
		Password:   hashPwd,
		Email:      email,
		Role:       role,
		Avatar:     Avatar,
		IP:         ip,
		Addr:       "内网地址", // 地址要根据IP算
		SignStatus: ctype.SignEmail,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
