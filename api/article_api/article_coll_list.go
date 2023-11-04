package article_api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/common"
	"gvb_server/utils/jwts"
)

// CollResponse 文章收藏界面需求
type CollResponse struct {
	models.ArticleModel
	CreatedAt string `json:"created_at"` // 收藏时间
}

func (ArticleApi) ArticleCollListView(c *gin.Context) {
	// 分页
	var cr models.PageInfo
	c.ShouldBindQuery(&cr)
	// 登录
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 分页查询
	var articleIDList []interface{}
	list, count, err := common.ComList(models.UserCollectModel{UserID: claims.UserID}, common.Option{
		PageInfo: cr,
	})

	// 收藏的文章id列表的获取,同时map映射对应文章id收藏的时间
	var collMap = map[string]string{}

	for _, model := range list {
		articleIDList = append(articleIDList, model.ArticleID)
		collMap[model.ArticleID] = model.CreatedAt.Format("2006-01-02 15:04:05")
	}

	// 传id列表，查es
	// 因为要查的是一个id列表，此处我们根据id查询的，且传入一个列表
	// 所以第一个参数要是一个string表示id，且表述为"_id"
	// 第二个参数要是一个数组，就是很多个参数的，进行查询，所以选择了
	// 也因为该查询函数的第二个参数类型为values ...interface{}，所以需要将articleIDList的类型修改成values ...interface{}
	boolSearch := elastic.NewTermsQuery("_id", articleIDList...)

	var collList = make([]CollResponse, 0)

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
		article.ID = hit.Id
		collList = append(collList, CollResponse{
			ArticleModel: article,
			CreatedAt:    collMap[hit.Id],
		})
	}
	res.OKWithList(collList, count, c)
}
