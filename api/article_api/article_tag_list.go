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

type TagsResponse struct {
	Tag           string   `json:"tag"`
	Count         int      `json:"count"`
	ArticleIDList []string `json:"article_id_list"` // 文章title列表
	CreatedAt     string   `json:"created_at"`
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

func (ArticleApi) ArticleTagListView(c *gin.Context) {
	// 分页查询
	var cr models.PageInfo
	_ = c.ShouldBindQuery(&cr)
	if cr.Limit == 0 {
		cr.Limit = 10
	}

	offset := (cr.Page - 1) * cr.Limit
	if offset < 0 {
		offset = 0
	}

	/*
		[{"tag": "python", "article_count": 2, "article_list": []}]
		es查询条件？
		elastic.NewBoolQuery() // 不要查询条件
	*/

	// 去重后标签的总数，这个总数单独查
	// elastic.NewCardinalityAggregation(). 会对指标进行去重
	// elastic.NewValueCountAggregation(). 不会进行去重操作
	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Aggregation("tags", elastic.NewCardinalityAggregation().Field("tags")).
		Size(0).
		Do(context.Background())
	cTag, _ := result.Aggregations.Cardinality("tags") // 标签总数
	count := int64(*cTag.Value)

	agg := elastic.NewTermsAggregation().Field("tags")
	// 在子集标签下的文章的列表中找到对应的keyword,此处keyword就是文章标题title
	agg.SubAggregation("article", elastic.NewTermsAggregation().Field("keyword"))
	// es查询的分页操作
	agg.SubAggregation("page", elastic.NewBucketSortAggregation().From(offset).Size(cr.Limit))

	query := elastic.NewBoolQuery()

	result, err = global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query).
		Aggregation("tags", agg).
		Size(0).
		Do(context.Background())
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage(err.Error(), c)
		return
	}

	var tagType TagsType
	var tagList = make([]*TagsResponse, 0) // 这个后面需要修改里面的createdAt所以需要用指针
	var tagStringList []string             // 标签title列表
	//fmt.Println(string(result.Aggregations["tags"]))
	_ = json.Unmarshal(result.Aggregations["tags"], &tagType)
	for _, bucket := range tagType.Buckets {
		var articleList []string
		for _, art := range bucket.Article.Buckets {
			articleList = append(articleList, art.Key)
		}

		tagList = append(tagList, &TagsResponse{
			Tag:           bucket.Key,
			Count:         bucket.DocCount,
			ArticleIDList: articleList,
		})
		tagStringList = append(tagStringList, bucket.Key) // 标签title列表
	}

	var tagModelList []models.TagModel
	global.DB.Find(&tagModelList, "title in ?", tagStringList)
	var tagDate = map[string]string{}
	//var tagDate map[string]string 是错的
	for _, model := range tagModelList {
		tagDate[model.Title] = model.CreatedAt.Format("2006-01-02 15:04:05")
	}

	for _, tagResponse := range tagList {
		tagResponse.CreatedAt = tagDate[tagResponse.Tag]
	}

	res.OKWithList(tagList, count, c)

}
