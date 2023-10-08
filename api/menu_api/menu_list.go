package menu_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

type Banner struct {
	ID   uint   `json:"id"`
	Path string `json:"path"`
}

type MenuResponse struct {
	models.MenuModel
	Banners []Banner `json:"banners"`
}

func (MenuApi) MenuListView(c *gin.Context) {
	// 先查菜单，得到菜单列表，同时查到所有菜单的ID列表
	var menuList []models.MenuModel
	var menuIDList []uint
	// sort desc 从大到小排
	global.DB.Order("sort desc").Find(&menuList).Select("id").Scan(&menuIDList)

	// 查连接表,找到对应menu_id关联的menu_banner
	var menuBanners []models.MenuBannerModel
	// .Preload("BannerModel")联表查询
	global.DB.Preload("BannerModel").Order("sort desc").Find(&menuBanners, "menu_id in ?", menuIDList)

	var menus []MenuResponse
	for _, model := range menuList {
		// model 就是一个菜单, 查找当前菜单对应的图片列表
		//var banners []Banner // 引用数据类型，如果声明了但是没有赋值等价于nil,前端得到的数据会显示null
		// 修改后的代码:
		var banners = make([]Banner, 0) // 分配一个长度为0的地址
		for _, banner := range menuBanners {
			if banner.MenuID == model.ID {
				banners = append(banners, Banner{
					ID:   banner.BannerID,
					Path: banner.BannerModel.Path, // 联表查询吼吼吼
				})
			}
		}
		menus = append(menus, MenuResponse{
			MenuModel: model,
			Banners:   banners,
		})
	}

	res.OKWithData(menus, c)
	return
}
