package advert_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/common"
	"strings"
)

// AdvertListView 广告列表分页查询
// @Tags 广告管理
// @summary 广告列表
// @Description 广告列表
// @Param data query models.PageInfo false "查询参数"
// @Router /api/adverts [get]
// @Produce json
// @success 200 {object} res.Response{data=res.ListResponse[models.AdvertModel]}
func (AdvertApi) AdvertListView(c *gin.Context) {
	var cr models.PageInfo
	if err := c.ShouldBindQuery(&cr); err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}

	referer := c.GetHeader("Referer")
	isShow := true
	if strings.Contains(referer, "admin") {
		// admin来的
		isShow = false
	}
	// 判断Referer是否包含admin，如果是（表示后台查询），就全部返回，不是就返回is_show=true的(就只展示需要展示的)
	list, count, _ := common.ComList(models.AdvertModel{IsShow: isShow}, common.Option{
		PageInfo: cr,
		Debug:    true,
	})

	res.OKWithList(list, count, c)
}
