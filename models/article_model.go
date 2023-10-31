package models

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gvb_server/global"
	"gvb_server/models/ctype"
)

// ArticleModel es存储文章，所以表结构也需要修改
type ArticleModel struct {
	ID        string `json:"id" structs:"id"`                 // es的id
	CreatedAt string `json:"created_at" structs:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at" structs:"updated_at"` // 更新时间

	Title    string `json:"title" structs:"title"`                // 文章标题
	KeyWord  string `json:"keyword,omit(list)" structs:"keyword"` // 关键字
	Abstract string `json:"abstract" structs:"abstract"`          // 文章简介
	Content  string `json:"content,omit(list)" structs:"content"` // 文章内容，在list的时候不要

	LookCount     int `json:"look_count" structs:"look_count"`         // 浏览量
	CommentCount  int `json:"comment_count" structs:"comment_count"`   // 评论量
	DiggCount     int `json:"digg_count" structs:"digg_count"`         // 点赞量
	CollectsCount int `json:"collects_count" structs:"collects_count"` // 收藏量

	UserID       uint   `json:"user_id" structs:"user_id"`               // 用户id
	UserNickName string `json:"user_nick_name" structs:"user_nick_name"` // 用户昵称
	UserAvatar   string `json:"user_avatar" structs:"user_avatar"`       // 用户头像

	Category string `json:"category" structs:"category"` // 文章分类
	Source   string `json:"source" structs:"source"`     // 文章来源
	Link     string `json:"link" structs:"link"`         // 原文连接

	// 后面便于不需要关联表查询，可以直接使用
	BannerID  uint   `json:"banner_id" structs:"banner_id"`   // 封面id
	BannerUrl string `json:"banner_url" structs:"banner_url"` // 文章封面

	Tags ctype.Array `json:"tags" structs:"tags"` // 文章标签
}

// Index 文章索引
func (ArticleModel) Index() string {
	return "article_index"
}

func (ArticleModel) Mapping() string {
	return `
	{
	  "settings": {
		"index":{
		  "max_result_window": "100000" 
		}
	  }, 
	  "mappings": {
		"properties": {
		  "title": { 
			"type": "text" 
		  },
		  "keyword": { 
			"type": "keyword" 
		  },
          "abstract": { 
			"type": "text" 
		  },
          "content": { 
			"type": "text" 
		  },
		  "look_count": {
			"type": "integer" 
		  },
		  "comment_count": {
			"type": "integer" 
		  },
		  "digg_count": {
			"type": "integer" 
		  },
		  "collects_count": {
			"type": "integer" 
		  },
		  "user_id": {
			"type": "integer" 
		  },
          "user_nick_name": { 
			"type": "keyword" 
		  },
          "user_avatar": { 
			"type": "keyword" 
		  },
          "category": { 
			"type": "keyword" 
		  },
          "source": { 
			"type": "keyword" 
		  },
          "link": { 
			"type": "keyword" 
		  },
		  "banner_id": {
			"type": "integer" 
		  },
          "banner_url": { 
			"type": "keyword" 
		  },
          "tags": { 
			"type": "keyword" 
		  },
		  "created_at":{
			"type": "date",
			"null_value": "null",
			"format": "yyyy-MM-dd HH:mm:ss"
		  },
		  "updated_at":{
			"type": "date",
			"null_value": "null",
			"format": "yyyy-MM-dd HH:mm:ss"
		  }
		}
	  }
	}
`
}

// IndexExists 判断索引是否存在
func (demo ArticleModel) IndexExists() bool {
	exists, err := global.ESClient.
		IndexExists(demo.Index()).
		Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return exists
	}
	return exists
}

// CreateIndex 创建索引
func (demo ArticleModel) CreateIndex() error {
	if demo.IndexExists() {
		// 有索引，就删掉索引，重新添加
		demo.RemoveIndex()
	}
	// 没有索引就创建索引
	createIndex, err := global.ESClient.
		CreateIndex(demo.Index()).
		BodyString(demo.Mapping()).
		Do(context.Background())

	if err != nil {
		logrus.Error("创建索引失败")
		fmt.Println(createIndex, err)
		logrus.Error(err.Error())
		return err
	}
	if !createIndex.Acknowledged {
		logrus.Error("创建索引失败")
		return err
	}
	logrus.Infof("索引 %s 创建成功", demo.Index())
	return nil
}

// RemoveIndex 删除索引
func (demo ArticleModel) RemoveIndex() error {
	logrus.Info("索引存在，删除索引")
	// 删除索引
	indexDelete, err := global.ESClient.DeleteIndex(demo.Index()).Do(context.Background())
	if err != nil {
		logrus.Error("删除索引失败")
		logrus.Error(err.Error())
		return err
	}

	if !indexDelete.Acknowledged {
		logrus.Error("删除索引失败")
		return err
	}
	logrus.Info("删除索引成功")
	return nil
}

// Create 添加文章的方法
func (demo ArticleModel) Create() (err error) {
	indexResponse, err := global.ESClient.Index().
		Index(demo.Index()).
		BodyJson(demo).Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	demo.ID = indexResponse.Id
	return nil
}

// ISExistData 是否存在该文章
// 这段代码是在检查一个名为ArticleModel的文章模型是否存在于Elasticsearch中。这是通过搜索具有特定标题的文章来完成的。
// 如果找到至少一篇文章，那么函数就会返回true，表示数据存在。
// 如果没有找到任何文章，或者在搜索过程中出现错误，函数就会返回false，表示数据不存在或无法确定。
func (demo ArticleModel) ISExistData() bool {
	res, err := global.ESClient.
		Search(demo.Index()).
		Query(elastic.NewTermQuery("keyword", demo.Title)).
		Size(1).
		Do(context.Background())
	// 这行代码使用全局Elasticsearch客户端执行搜索查询。
	// 它在特定索引（由demo.Index()指定）中搜索包含特定标题（由demo.Title指定）的文章。
	// 它只请求一个结果（.Size(1)），然后执行查询（.Do(context.Background())）。
	if err != nil {
		logrus.Error(err.Error())
		return false
	}
	// 如果查询结果的总数大于0，这意味着找到了至少一篇匹配的文章，所以返回true。
	if res.Hits.TotalHits.Value > 0 {
		return true
	}
	return false
}
