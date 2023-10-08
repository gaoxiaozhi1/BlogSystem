package menu_api

import (
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

func (MenuApi) MenuUpdateView(c *gin.Context) {
	var cr MenuRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	id := c.Param("id")

	// 先把之前的banner清空
	var menuModel models.MenuModel
	err = global.DB.Take(&menuModel, id).Error
	if err != nil {
		res.FailWithMessage("菜单不存在", c)
		return
	}
	// 联表，把第三张表清空
	global.DB.Model(&menuModel).Association("Banners").Clear()

	// 如果选择了banner,那就添加
	if len(cr.ImageSortList) > 0 {
		// 操作第三张表
		var menuBannerList []models.MenuBannerModel
		for _, imageSort := range cr.ImageSortList {
			menuBannerList = append(menuBannerList, models.MenuBannerModel{
				MenuID:   menuModel.ID,
				BannerID: imageSort.ImageID,
				Sort:     imageSort.Sort,
			})
		}
		err = global.DB.Create(&menuBannerList).Error
		if err != nil {
			res.FailWithMessage("创建菜单图片失败", c)
			return
		}
	}

	// 普通更新，
	// 因为要防范零止问题，所以转化成map来更新数据
	// 结构体转map的第三方包structs
	maps := structs.Map(&cr)
	// Updates 好用嗷嗷嗷嗷嗷嗷嗷嗷嗷嗷
	err = global.DB.Model(&menuModel).Updates(maps).Error

	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("修改菜单失败", c)
		return
	}

	res.OKWithMessage("修改菜单成功", c)
}
