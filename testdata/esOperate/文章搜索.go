package main

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/models"
)

func main() {
	// 读取配置文件
	core.InitConf()
	// 初始化日志
	global.Log = core.InitLogger()
	// 连接es
	global.ESClient = core.EsConnect()

	boolSearch := elastic.NewBoolQuery()
	boolSearch.Must(
		//elastic.NewMatchQuery("title", key),
		// 搜索多个，如果搜到了，那么这些部分就会返回高亮
		elastic.NewMultiMatchQuery("title", "abstract", "content"),
	)
	// es查询得到最终展示部分
	res, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(boolSearch).
		Highlight(elastic.NewHighlight().Field("title")). // 标题高亮
		Size(100).
		Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	for _, hit := range res.Hits.Hits {
		fmt.Println(string(hit.Source))
		fmt.Println(hit.Highlight) // 高光
	}
}
