package tag_api

import (
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

func (TagApi) TagUpdateView(c *gin.Context) {
	id := c.Param("id")
	var cr TagRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}

	// 判断广告是否存在
	var tag models.TagModel
	err = global.DB.Take(&tag, id).Error
	if err != nil { // 如果数据库中存在
		res.FailWithMessage("标签不存在", c)
		return
	}

	// 结构体转map的第三方包structs
	maps := structs.Map(&cr)
	// Updates 好用嗷嗷嗷嗷嗷嗷嗷嗷嗷嗷
	err = global.DB.Model(&tag).Updates(maps).Error

	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("修改标签失败", c)
		return
	}

	res.OKWithMessage("修改标签成功", c)
}
