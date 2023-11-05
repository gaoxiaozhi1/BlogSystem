package article_api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
)

// FullTextSearchView 全文搜索，不需要做分页， 根据关键字查询
func (ArticleApi) FullTextSearchView(c *gin.Context) {
	var cr models.PageInfo
	_ = c.ShouldBindQuery(&cr)

	// 不搜的时候也能搜到
	boolSearch := elastic.NewBoolQuery()
	if cr.Key != "" {
		boolSearch.Must(elastic.NewMultiMatchQuery(cr.Key, "title", "body"))
	}
	// es查询得到最终展示部分
	result, err := global.ESClient.
		Search(models.FullTextModel{}.Index()).
		Query(boolSearch).
		Highlight(elastic.NewHighlight().Field("body")). // 文章内容中关键字高亮
		Size(100).
		Do(context.Background())
	if err != nil {
		return
	}

	// 搜索到的总条数
	count := result.Hits.TotalHits.Value
	fullTextList := make([]models.FullTextModel, 0)
	for _, hit := range result.Hits.Hits {
		var model models.FullTextModel
		json.Unmarshal(hit.Source, &model)
		body, ok := hit.Highlight["body"]
		if ok {
			model.Body = body[0]
		}
		fullTextList = append(fullTextList, model)
	}
	res.OKWithList(fullTextList, count, c)
}
