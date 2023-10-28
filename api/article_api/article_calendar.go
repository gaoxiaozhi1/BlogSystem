package article_api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"time"
)

type CalendarResponse struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type BucketsType struct {
	Buckets []struct {
		KeyAsString string `json:"key_as_string"`
		Key         int64  `json:"key"`
		DocCount    int    `json:"doc_count"`
	} `json:"buckets"`
}

var DateCount = map[string]int{}

// ArticleCalendarView 发布文章的日历显示
func (ArticleApi) ArticleCalendarView(c *gin.Context) {
	// 时间聚合，日历按天聚合
	// 这行代码创建了一个新的日期直方图聚合，该聚合根据文章的创建日期（"created_at"字段）进行聚合，并且聚合间隔为一天。
	agg := elastic.NewDateHistogramAggregation().Field("created_at").CalendarInterval("day")

	// 时间段搜索
	// 从今天开始，到去年的今天
	now := time.Now()
	// 这些代码获取当前时间，计算一年前的日期，并定义一个时间格式字符串。
	aYearAgo := now.AddDate(-1, 0, 0) // 一年前
	//fmt.Println(now, aYearAgo)
	//twoHourAgo := now.Add(-2 * time.Hour) // 两小时前
	format := "2006-01-02 15:04:05" // 时间格式
	// 这行代码创建了一个范围查询，该查询查找创建日期在一年前和现在之间的文章。
	// lt 小于 gt 大于 （此处为以天为单位的查询区间为一年，大于一年前，小于今天的日期之间）
	query := elastic.NewRangeQuery("created_at").
		Gte(aYearAgo.Format(format)).
		Lte(now.Format(format))

	// 这行代码使用全局Elasticsearch客户端执行搜索查询。
	// 它在特定索引（由models.ArticleModel{}.Index()指定）中执行查询，并添加了之前定义的聚合。
	// 它不请求任何文档（.Size(0)），然后执行查询（.Do(context.Background())）。
	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query). // 查询的区间为一年
		Aggregation("calendar", agg).
		Size(0).
		Do(context.Background())
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("查询失败", c)
		return
	}

	var data BucketsType
	_ = json.Unmarshal(result.Aggregations["calendar"], &data)

	var resList = make([]CalendarResponse, 0) // 初始化，因为可能为空值

	for _, bucket := range data.Buckets {
		// 这行代码将当前桶的键（即时间）从字符串转换为time.Time类型。它假设键是按照某种格式（由format变量指定）编码的。
		Time, _ := time.Parse(format, bucket.KeyAsString)
		DateCount[Time.Format("2006-01-02")] = bucket.DocCount
	}

	days := int(now.Sub(aYearAgo).Hours() / 24) // 去年这天到今年今天总共有多少天
	for i := 0; i <= days; i++ {
		day := aYearAgo.AddDate(0, 0, i).Format("2006-01-02") // 从去年那天开始按天累加，格式为一天一天

		count, _ := DateCount[day]
		resList = append(resList, CalendarResponse{
			Date:  day,
			Count: count,
		})
	}

	res.OKWithData(resList, c)
}
