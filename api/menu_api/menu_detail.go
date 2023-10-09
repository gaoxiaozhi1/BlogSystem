package menu_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

// MenuDetailView 查单个菜单的细节嗷
func (MenuApi) MenuDetailView(c *gin.Context) {
	// 根据id查菜单
	id := c.Param("id")
	var menuModel models.MenuModel
	err := global.DB.Take(&menuModel, id).Error
	if err != nil {
		res.FailWithMessage("菜单不存在", c)
		return
	}

	// 查连接表
	var menuBannerList []models.MenuBannerModel
	global.DB.Preload("BannerModel").Order("sort desc").Find(&menuBannerList, "menu_id = ?", id)

	// 查对应菜单的图片表
	var bannerList = make([]Banner, 0)
	for _, banner := range menuBannerList {
		bannerList = append(bannerList, Banner{
			ID:   banner.BannerID,
			Path: banner.BannerModel.Path,
		})
	}
	menuResponse := MenuResponse{
		MenuModel: menuModel,
		Banners:   bannerList,
	}
	res.OKWithData(menuResponse, c)
	return
}
