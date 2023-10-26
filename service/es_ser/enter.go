package es_ser

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gvb_server/global"
	"gvb_server/models"
)

// CommList 公共的查询列表页，文章列表
// 文章列表的后台部分和前端展示页，共同的就是都不显示文章内容，只会显示文章简介
// 所以此处需要过滤掉文章内容的显示
func CommList(key string, page, limit int) (list []models.ArticleModel, count int, err error) {
	boolSearch := elastic.NewBoolQuery()
	from := page
	if key != "" {
		boolSearch.Must(
			elastic.NewMatchQuery("title", key),
		)
	}
	if limit == 0 {
		limit = 10
	}
	if from == 0 {
		from = 1
	}

	// es查询得到最终展示部分
	res, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(boolSearch).
		// 排除掉文章内容部分， 但是这种方法比较死板
		// 可以调用一个库 github.com/liu-cn/json-filter/filter
		//FetchSourceContext(
		//	elastic.NewFetchSourceContext(true).Exclude("content"),
		//).
		From((from - 1) * limit).
		Size(limit).
		Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	count = int(res.Hits.TotalHits.Value) // 搜索到结果的总条数
	demoList := []models.ArticleModel{}
	for _, hit := range res.Hits.Hits {
		var model models.ArticleModel
		data, err := hit.Source.MarshalJSON() // 类型转换
		if err != nil {
			logrus.Error(err.Error())
			continue
		}
		err = json.Unmarshal(data, &model) // json->model类型
		if err != nil {
			logrus.Error(err.Error())
			continue
		}
		model.ID = hit.Id
		demoList = append(demoList, model)
	}
	return demoList, count, err
}
