package menu_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

func (MenuApi) MenuRemoveView(c *gin.Context) {
	var cr models.RemoveRequest // 要删除的广告的id列表
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}

	// 先从数据库中查询对应的信息
	var menutList []models.MenuModel
	count := global.DB.Find(&menutList, cr.IDList).RowsAffected
	if count == 0 {
		res.OKWithMessage("菜单不存在", c)
		return
	}

	// 查到啦，就要批量删除喽
	// 要先删除关联的第三张表，然后删除菜单表----因为要删除两个东西，所以用事务
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 删除关联的第三张表
		err = global.DB.Model(&menutList).Association("Banners").Clear()
		if err != nil {
			global.Log.Error(err)
			return err
		}

		// 删除菜单表
		global.DB.Delete(&menutList)
		if err != nil {
			global.Log.Error(err)
			return err
		}

		// 都成功
		return nil
	})

	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("删除菜单失败", c)
		return
	}
	res.OKWithMessage(fmt.Sprintf("共删除 %d 个菜单", count), c)

}
