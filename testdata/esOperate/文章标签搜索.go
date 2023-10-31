package main

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/models"
)

type TagsResponse struct {
	Tag           string   `json:"tag"`
	Count         int      `json:"count"`
	ArticleIDList []string `json:"article_id_list"` // 文章title列表
}

type TagsType struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
		Article  struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"article"`
	} `json:"buckets"`
}

func main() {
	// 读取配置文件
	core.InitConf()
	// 初始化日志
	global.Log = core.InitLogger()
	// 连接es
	global.ESClient = core.EsConnect()

	/*
		[{"tag": "python", "article_count": 2, "article_list": []}]
		es查询条件？
		elastic.NewBoolQuery() // 不要查询条件
	*/

	agg := elastic.NewTermsAggregation().Field("tags")
	// 在子集标签下的文章的列表中找到对应的keyword,此处keyword就是文章标题title
	agg.SubAggregation("article", elastic.NewTermsAggregation().Field("keyword"))
	//agg.SubAggregation("article_key", elastic.NewTermsAggregation().Field("keyword"))
	// 同样可以查询文章id
	//agg.SubAggregation("article_id", elastic.NewTermsAggregation().Field("_id"))
	// 如果按上面两种写，就可以两种分开显示出来，那么接收显示的结构体同样需要修改。不能两个Field放在一起显示出来会冲突的。。。

	query := elastic.NewBoolQuery()

	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query).
		Aggregation("tags", agg).
		Size(0).
		Do(context.Background())
	if err != nil {
		global.Log.Error(err)
		return
	}

	var tagType TagsType
	var tagList = make([]TagsResponse, 0)
	_ = json.Unmarshal(result.Aggregations["tags"], &tagType)
	for _, bucket := range tagType.Buckets {
		var articleList []string
		for _, art := range bucket.Article.Buckets {
			articleList = append(articleList, art.Key)
		}

		tagList = append(tagList, TagsResponse{
			Tag:           bucket.Key,
			Count:         bucket.DocCount,
			ArticleIDList: articleList,
		})
	}

	fmt.Println(tagList)

}
