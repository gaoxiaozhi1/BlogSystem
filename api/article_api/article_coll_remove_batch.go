package article_api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/es_ser"
	"gvb_server/utils/jwts"
)

func (ArticleApi) ArticleCollBatchRemoveView(c *gin.Context) {
	var cr models.ESIDListRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}
	// 参数绑定
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	var collects []models.UserCollectModel
	var articleIDList []string
	// 先Find再Selete，如果先Selete，那么Find就是Selete中筛选的东西啦，就搜不到
	global.DB.Find(&collects, "user_id = ? and article_id in ?", claims.UserID, cr.IDList).
		Select("article_id").
		Scan(&articleIDList)
	if len(articleIDList) == 0 {
		res.FailWithMessage("请求非法", c)
		return
	}
	var artIDList []interface{}
	for _, artID := range articleIDList {
		artIDList = append(artIDList, artID)
	}
	// 更新文章数
	boolSearch := elastic.NewTermsQuery("_id", artIDList...)
	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(boolSearch).
		Size(1000).
		Do(context.Background())
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}

	for _, hit := range result.Hits.Hits {
		var article models.ArticleModel
		err = json.Unmarshal(hit.Source, &article)
		if err != nil {
			global.Log.Error(err)
			continue
		}

		// 更新文章收藏操作
		err = es_ser.ArticleUpdate(hit.Id, map[string]any{
			"collects_count": article.CollectsCount - 1,
		})
		if err != nil {
			global.Log.Error(err)
			continue
		}
		// 删除
	}
	// 删除用户收藏表对应内容
	global.DB.Delete(&collects)
	res.OKWithMessage(fmt.Sprintf("成功取消收藏 %d 篇文章", len(articleIDList)), c)
}
