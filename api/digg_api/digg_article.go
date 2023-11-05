package digg_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/redis_ser"
)

// DiggArticleView 文章点赞操作
func (DiggApi) DiggArticleView(c *gin.Context) {
	var cr models.ESIDRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}
	// 对长度校验，
	// 查es
	redis_ser.NewDigg().Set(cr.ID)
	res.OKWithMessage("文章点赞成功", c)
}
