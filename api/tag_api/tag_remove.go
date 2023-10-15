package tag_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

func (TagApi) TagRemoveView(c *gin.Context) {
	var cr models.RemoveRequest // 要删除的广告的id列表
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}

	// 先从数据库中查询对应的信息
	var tagList []models.TagModel
	count := global.DB.Find(&tagList, cr.IDList).RowsAffected
	if count == 0 {
		res.OKWithMessage("标签不存在", c)
		return
	}

	// 查到啦，就要批量删除喽
	// 如果这个标签下有文章怎么办
	var articleCount int64
	global.DB.Model(&models.ArticleModel{}).Where("tag_id IN ?", cr.IDList).Count(&articleCount)

	if articleCount > 0 {
		res.FailWithMessage("该标签下存在文章，无法删除", c)
		return
	}

	// 标签下没有文章，可以进行批量删除
	global.DB.Delete(&tagList)
	res.OKWithMessage(fmt.Sprintf("共删除 %d 个标签", count), c)

}
