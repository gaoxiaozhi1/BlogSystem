package images_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

type ImageUpdateRequest struct {
	ID   uint   `json:"id" binding:"required" msg:"请选择文件id"`
	Name string `json:"name" binding:"required" msg:"请选择文件名称"` // 新名字
}

// ImageUpdateView 修改图片名称(没有修改存储路径上图片的名字，只改数据库中的name)
func (ImagesApi) ImageUpdateView(c *gin.Context) {
	// 转换前端传来的数据
	var cr ImageUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		// 根据报错信息，确定是因为什么数据没传导致的报错
		res.FailWithError(err, &cr, c)
		return
	}

	// 查找是否存在
	var imageModel models.BannerModel
	err = global.DB.Take(&imageModel, cr.ID).Error
	if err != nil {
		res.FailWithMessage("文件不存在", c)
		return
	}
	// 存在就修改
	err = global.DB.Model(&imageModel).Update("name", cr.Name).Error
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}

	res.OKWithMessage("图片名称修改成功", c)
	return
}
