package main

import (
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/service/redis_ser"
)

func main() {
	core.InitConf()
	global.Log = core.InitLogger()
	global.Redis = core.ConnectRedis()
	global.ESClient = core.EsConnect()

	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).
		Query(elastic.NewMatchAllQuery()).
		Size(10000).
		Do(context.Background())

	if err != nil {
		logrus.Error(err)
		return
	}

	diggInfo := redis_ser.NewDigg().GetInfo()
	lookInfo := redis_ser.NewArticleLook().GetInfo()
	for _, hit := range result.Hits.Hits {
		var article models.ArticleModel
		err = json.Unmarshal(hit.Source, &article)
		digg := diggInfo[hit.Id]
		look := lookInfo[hit.Id]
		if digg == 0 && look == 0 {
			logrus.Info(article.Title, "点赞数和浏览数无变化")
			continue
		}
		article.DiggCount += digg
		article.LookCount += look

		//fmt.Println(article.Title, hit.Id, article.DiggCount)
		_, err := global.ESClient.
			Update().
			Index(models.ArticleModel{}.Index()).
			Id(hit.Id).
			Doc(map[string]int{
				"digg_count": article.DiggCount,
				"look_count": article.LookCount,
			}).
			Do(context.Background())
		if err != nil {
			logrus.Error(err.Error())
			continue
		}
		logrus.Infof("%s, 点赞数同步成功，点赞数: %d 浏览量 %d", article.Title, article.DiggCount, article.LookCount)
	}
	// 清空redis中的数据
	redis_ser.NewDigg().Clear()
	redis_ser.NewArticleLook().Clear()
}
