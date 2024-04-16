package article_api

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/ctype"
	"gvb_server/models/res"
	"gvb_server/service/es_ser"
	"time"
)

type ArticleUpdateRequest struct {
	Title    string   `json:"title"`     // 文章标题
	Abstract string   `json:"abstract"`  // 文章简介
	Content  string   `json:"content"`   // 文章内容
	Category string   `json:"category"`  // 文章分类
	Source   string   `json:"source"`    // 文章来源
	Link     string   `json:"link"`      // 原文连接
	BannerID uint     `json:"banner_id"` // 文章封面id
	Tags     []string `json:"tags"`      // 文章标签
	ID       string   `json:"id"`
}

func (ArticleApi) ArticleUpdateView(c *gin.Context) {
	// 参数绑定
	var cr ArticleUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		global.Log.Error(err)
		res.FailWithError(err, &cr, c)
		return
	}

	var bannerUrl string
	if cr.BannerID != 0 {
		err = global.DB.Model(models.BannerModel{}).Where("id = ?", cr.BannerID).Select("path").Scan(&bannerUrl).Error
		if err != nil {
			res.FailWithMessage("banner不存在", c)
			return
		}
	}

	article := models.ArticleModel{
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Title:     cr.Title,
		KeyWord:   cr.Title,
		Abstract:  cr.Abstract,
		Content:   cr.Content,
		Category:  cr.Category,
		Source:    cr.Source,
		Link:      cr.Link,
		BannerID:  cr.BannerID,
		BannerUrl: bannerUrl,
		Tags:      cr.Tags,
	}

	// 结构体转map
	maps := structs.Map(&article)
	var DataMap = map[string]any{}
	// 去掉map中的空值
	for k, v := range maps {
		switch val := v.(type) {
		case string:
			if val == "" {
				continue
			}
		case uint:
			if val == 0 {
				continue
			}
		case int:
			if val == 0 {
				continue
			}
		case ctype.Array:
			if len(val) == 0 {
				continue
			}
		case []string:
			if len(val) == 0 {
				continue
			}
		}
		DataMap[k] = v
	}

	//fmt.Println(DataMap)
	// 判断文章是否存在， 这里之后article会变成之前的article的值
	err = article.GetDataByID(cr.ID)
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("文章不存在", c)
		return
	}
	//fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	//fmt.Println(article)

	err = es_ser.ArticleUpdate(cr.ID, DataMap)
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("文章更新失败", c)
		return
	}

	// 更新成功，同步数据到全文搜索
	newArticle, _ := es_ser.CommDetail(cr.ID)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println(newArticle)
	if article.Content != newArticle.Content || article.Title != newArticle.Title {
		es_ser.DeleteFullTextByArticleID(cr.ID)
		es_ser.AsyncArticleByFullText(cr.ID, newArticle.Title, newArticle.Content)
	}

	res.OKWithMessage("更新成功", c)
}
