package article_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/es_ser"
	"gvb_server/service/redis_ser"
)

func (ArticleApi) ArticleDetailView(c *gin.Context) {
	var cr models.ESIDRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}
	// 点进文章详情页，则浏览量++
	redis_ser.Look(cr.ID)
	
	model, err := es_ser.CommDetail(cr.ID)
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}
	res.OKWithData(model, c)
}

type ArticleDetailRequest struct {
	Title string `json:"title" form:"title"`
}

func (ArticleApi) ArticleDetailByTitleView(c *gin.Context) {
	var cr ArticleDetailRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}
	model, err := es_ser.CommDetailByKeyWord(cr.Title)
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}
	res.OKWithData(model, c)

}
