package es_ser

import (
	"context"
	"encoding/json"
	"errors"
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

// CommDetail 文章详情页根据用户id查找
func CommDetail(id string) (model models.ArticleModel, err error) {
	// 查询详情
	res, err := global.ESClient.Get().
		Index(models.ArticleModel{}.Index()).
		Id(id).
		Do(context.Background())
	if err != nil {
		return
	}
	err = json.Unmarshal(res.Source, &model)
	if err != nil {
		return
	}
	model.ID = res.Id
	return
}

// CommDetailByKeyWord 这段代码定义了一个名为CommDetailByKeyWord的函数，
// 该函数接收一个关键字作为参数，并返回一个ArticleModel对象和一个错误。
// 这个函数的主要目的是在Elasticsearch中查找具有特定关键字的文章，并返回找到的第一篇文章的详细信息。
func CommDetailByKeyWord(key string) (model models.ArticleModel, err error) {
	// 查询详情
	res, err := global.ESClient.Search().
		Index(models.ArticleModel{}.Index()).
		Query(elastic.NewTermQuery("keyword", key)).
		Size(1).
		Do(context.Background())
	// 这行代码使用全局Elasticsearch客户端执行搜索查询。
	// 它在特定索引（由models.ArticleModel{}.Index()指定）中搜索包含特定关键字（由key指定）的文章。
	// 它只请求一个结果（.Size(1)），然后执行查询（.Do(context.Background())）。
	if err != nil {
		return
	}

	if res.Hits.TotalHits.Value == 0 {
		return model, errors.New("文章不存在")
	}
	// 这行代码获取搜索结果中的第一篇文章。
	hit := res.Hits.Hits[0]
	// 这行代码将搜索结果中的第一篇文章的源数据解析为ArticleModel对象。
	err = json.Unmarshal(hit.Source, &model)
	if err != nil {
		return
	}
	model.ID = hit.Id
	// 在Elasticsearch中，每个文档都有一个唯一的ID，这个ID用于标识和检索文档。在这段代码中，hit.Id就是搜索结果中第一篇文章的ID。
	// model.ID = hit.Id这行代码的作用是将这个ID保存到ArticleModel对象中。
	// 这样做的好处是，以后如果需要对这篇文章进行更新或删除操作，就可以直接使用model.ID来找到这篇文章。
	// 总的来说，这行代码是在将Elasticsearch中的文档ID与我们的ArticleModel对象关联起来，以便于后续的数据操作。
	return
}
