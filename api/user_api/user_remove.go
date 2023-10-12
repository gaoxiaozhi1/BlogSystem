package user_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

func (UserApi) UserRemoveView(c *gin.Context) {
	var cr models.RemoveRequest // 要删除的用户的id列表
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}

	// 先从数据库中查询对应的信息
	var userList []models.UserModel
	count := global.DB.Find(&userList, cr.IDList).RowsAffected
	if count == 0 {
		res.OKWithMessage("用户不存在", c)
		return
	}

	// 查到啦，就要批量删除喽
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// TODO:别除用户，消息表，评论表，用户收藏的文章，用户发布的文章
		// 删除用户
		global.DB.Delete(&userList)
		if err != nil {
			global.Log.Error(err)
			return err
		}
		// 都成功
		return nil
	})

	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("删除用户失败", c)
		return
	}
	res.OKWithMessage(fmt.Sprintf("共删除 %d 个用户", count), c)

}
