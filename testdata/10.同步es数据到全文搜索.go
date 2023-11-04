package main

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/service/es_ser"
)

func main() {
	// 读取配置文件
	core.InitConf()
	// 初始化日志
	global.Log = core.InitLogger()
	// 连接es
	global.ESClient = core.EsConnect()

	// 先查
	boolSearch := elastic.NewMatchAllQuery()
	res, _ := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(boolSearch).
		Size(1000).
		Do(context.Background())

	for _, hit := range res.Hits.Hits {
		var article models.ArticleModel
		_ = json.Unmarshal(hit.Source, &article)

		// 批量添加
		indexList := es_ser.GetSearchIndexByContent(hit.Id, article.Title, article.Content)

		// 先创建一个桶
		bulk := global.ESClient.Bulk()
		for _, indexData := range indexList {
			req := elastic.NewBulkIndexRequest().Index(models.FullTextModel{}.Index()).Doc(indexData) // 即将加入桶的请求
			bulk.Add(req)
		}
		result, err := bulk.Do(context.Background()) // 执行，统一删除
		if err != nil {
			logrus.Error(err)
			continue
		}
		fmt.Println(article.Title, "添加成功, 共", len(result.Succeeded()), "条!")

	}
}
