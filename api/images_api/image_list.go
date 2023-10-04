package images_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/common"
)

// ImageListView 图片列表查询页
func (ImagesApi) ImageListView(c *gin.Context) {
	// 分页
	var cr models.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}

	list, count, err := common.ComList(models.BannerModel{}, common.Option{
		PageInfo: cr,
		Debug:    true,
	})

	res.OKWithList(list, count, c)
	return
}
