package advert_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

// AdvertRemoveView 批量删除广告，根据idList从数据库中进行批量删除
// @Tags 广告管理
// @summary 批量删除广告
// @Description 批量删除广告
// @Param data body models.RemoveRequest true "广告的id列表"
// @Router /api/adverts [delete]
// @Produce json
// @success 200 {object} res.Response{data=string}
func (AdvertApi) AdvertRemoveView(c *gin.Context) {
	var cr models.RemoveRequest // 要删除的广告的id列表
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}

	// 先从数据库中查询对应的信息
	var advertList []models.AdvertModel
	count := global.DB.Find(&advertList, cr.IDList).RowsAffected
	if count == 0 {
		res.OKWithMessage("广告不存在", c)
		return
	}

	// 查到啦，就要批量删除喽
	global.DB.Delete(&advertList)
	res.OKWithMessage(fmt.Sprintf("共删除 %d 个广告", count), c)

}
