package article_api

import (
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/es_ser"
)

// ArticleSearchRequest 因为有根据标签搜索的，这就需要传参，但是分页model中没有tag的，所以需要重写一个
type ArticleSearchRequest struct {
	models.PageInfo
	Tag string `json:"tag" form:"tag"` // 查询是的uri中叫tag而不是tags
}

// ArticleListView 文章列表的后台部分和前端展示页，共同的就是都不显示文章内容，只会显示文章简介
func (ArticleApi) ArticleListView(c *gin.Context) {
	var cr ArticleSearchRequest
	if err := c.ShouldBindQuery(&cr); err != nil {
		res.FailWitheCode(res.ArgumentError, c) // 参数错误
		return
	}
	//fmt.Println(cr.Tag)
	list, count, err := es_ser.CommList(es_ser.Option{
		PageInfo: cr.PageInfo,
		Fields:   []string{"title", "content", "abstract"},
		Tag:      cr.Tag,
	})
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("查询失败", c)
		return
	}

	// json-filter空值问题
	data := filter.Omit("list", list)
	_list, _ := data.(filter.Filter)
	if _list.MustJSON() == "{}" {
		data = make([]models.ArticleModel, 0)
	}

	res.OKWithList(data, int64(count), c)
}
