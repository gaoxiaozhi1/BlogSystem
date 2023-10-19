package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gvb_server/core"
)

var client *elastic.Client

func EsConnect() *elastic.Client {
	var err error
	sniffOpt := elastic.SetSniff(false)

	host := "http://127.0.0.1:9200"

	client, err = elastic.NewClient( // 可以设置超时时间
		elastic.SetURL(host),
		sniffOpt,
		elastic.SetBasicAuth("", ""),
	)

	if err != nil {
		logrus.Fatalf("es连接失败%s", err.Error())
	}
	return client
}

// init 自动执行
func init() {
	// 加载数据库配置
	core.InitConf()
	core.InitLogger()
	client = EsConnect()
}

type DemoModel struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	UserID   uint   `json:"user_id"`
	CreateAt string `json:"created_at"`
}

func (DemoModel) Index() string {
	return "demo_index"
}

// Create 创建文章
// data *DemoModel指针便于修改data.ID
func Create(data *DemoModel) (err error) {
	indexResponse, err := client.Index().
		Index(data.Index()).
		// 转成json
		BodyJson(data).Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	data.ID = indexResponse.Id
	return nil

}

// FindList 列表查询 ,返回全部的文章信息
func FindList(key string, page, limit int) (demoList []DemoModel, count int) {
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

	res, err := client.
		Search(DemoModel{}.Index()).
		Query(boolSearch).
		From((from - 1) * limit).
		Size(limit).
		Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	count = int(res.Hits.TotalHits.Value) //搜索到结果总条数
	for _, hit := range res.Hits.Hits {
		var demo DemoModel
		data, err := hit.Source.MarshalJSON()
		if err != nil {
			logrus.Error(err.Error())
			continue
		}
		err = json.Unmarshal(data, &demo)
		if err != nil {
			logrus.Error(err)
			continue
		}
		demo.ID = hit.Id
		demoList = append(demoList, demo)
	}
	return demoList, count
}

// FindSourceList 搜索，检索那类
func FindSourceList(key string, page, limit int) {
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

	res, err := client.
		Search(DemoModel{}.Index()).
		Query(boolSearch).
		Source(`{"_source": ["title"]}`).
		From((from - 1) * limit).
		Size(limit).
		Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	count := int(res.Hits.TotalHits.Value) //搜索到结果总条数
	demoList := []DemoModel{}
	for _, hit := range res.Hits.Hits {
		var demo DemoModel
		data, err := hit.Source.MarshalJSON()
		if err != nil {
			logrus.Error(err.Error())
			continue
		}
		err = json.Unmarshal(data, &demo)
		if err != nil {
			logrus.Error(err)
			continue
		}
		demo.ID = hit.Id
		demoList = append(demoList, demo)
	}
	fmt.Println(demoList, count)
}

func Update(id string, data *DemoModel) error {
	_, err := client.
		Update().
		Index(DemoModel{}.Index()).
		Id(id).
		Doc(map[string]string{
			"title": data.Title,
		}).
		Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	logrus.Info("更新demo成功")
	return nil
}

// Remove 批量删除，根据id列表，返回删除个数
func Remove(idList []string) (count int, err error) {
	// 弄一个桶
	bulkService := client.Bulk().Index(DemoModel{}.Index()).Refresh("true") // 更新
	for _, id := range idList {
		req := elastic.NewBulkDeleteRequest().Id(id) // 把要删除的文章id加入到桶里
		bulkService.Add(req)
	}
	res, err := bulkService.Do(context.Background()) // 执行，统一删除
	return len(res.Succeeded()), err
}

func main() {
	//DemoModel{}.CreateIndex()
	//Create(&DemoModel{Title: "python", UserID: 2, CreateAt: time.Now().Format("2006-01-02 15:04:05")})
	//list, count := FindList("", 1, 10)
	//fmt.Println(list, count)
	//FindSourceList("python", 1, 10) // 搜索似乎失效了
	//Update("H2GGRosBofK7yM7TUrX6", &DemoModel{Title: "Go语言开发基础"})
	//count, err := Remove([]string{"IGGHRosBofK7yM7TbrVJ"})
	//fmt.Println(count, err)
}
