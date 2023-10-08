package images_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

type ImageResponse struct {
	ID   uint   `json:"id"`
	Path string `json:"path"` // 路径图片
	Name string `json:"name"` // 图片名称

}

// ImageNameListView 图片列表查询页
// @Tags 图片管理
// @summary 图片名称列表
// @Description 图片名称列表
// @Router /api/image_names [get]
// @Produce json
// @success 200 {object} res.Response{data=[]ImageResponse}
func (ImagesApi) ImageNameListView(c *gin.Context) {
	var imageList []ImageResponse
	global.DB.Model(&models.BannerModel{}).Select("id", "path", "name").Scan(&imageList)
	res.OKWithData(imageList, c)
}
