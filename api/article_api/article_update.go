package article_api

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/ctype"
	"gvb_server/models/res"
	"time"
)

type ArticleUpdateRequest struct {
	Title    string `json:"title"`    // 文章标题
	Abstract string `json:"abstract"` // 文章简介
	Content  string `json:"content"`  // 文章内容
	Category string `json:"category"` // 文章分类
	Source   string `json:"source"`   // 文章来源
	Link     string `json:"link"`     // 原文连接
	// 后面便于不需要关联表查询，可以直接使用
	BannerID uint        `json:"banner_id"` // 文章封面id
	Tags     ctype.Array `json:"tags"`      // 文章标签
	ID       string      `json:"id"`
}

func (ArticleApi) ArticleUpdateView(c *gin.Context) {
	// 参数绑定
	var cr ArticleUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		global.Log.Error(err.Error())
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
	fmt.Println(DataMap)
	_, err = global.ESClient.
		Update().
		Index(models.ArticleModel{}.Index()).
		Id(cr.ID).
		Doc(DataMap).
		Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		res.FailWithMessage("更新失败", c)
		return
	}
	res.OKWithMessage("更新成功", c)
}
