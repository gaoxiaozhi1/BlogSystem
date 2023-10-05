package images_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

// 批量删除图片，根据idList从数据库中进行批量删除
func (ImagesApi) ImageRemoveView(c *gin.Context) {
	var cr models.RemoveRequest // 要删除的图片的id列表
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}

	// 先从数据库中查询对应的信息
	var imageList []models.BannerModel
	count := global.DB.Find(&imageList, cr.IDList).RowsAffected
	if count == 0 {
		res.OKWithMessage("文件不存在", c)
		return
	}

	// 查到啦，就要批量删除喽
	global.DB.Delete(&imageList)
	res.OKWithMessage(fmt.Sprintf("共删除 %d 张图片", count), c)

}
