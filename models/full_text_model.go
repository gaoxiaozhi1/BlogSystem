package models

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gvb_server/global"
)

// FullTextModel 全文搜索的索引，也是es的表
type FullTextModel struct {
	ID    string `json:"id" structs:"id"`       // es的id
	Key   string `json:"key"`                   // 文章关联的id，这样才能在根据文章title删除文章时删除全文搜素的文章
	Title string `json:"title" structs:"title"` // 标题
	Slug  string `json:"slug" structs:"slug"`   // 跳转地址（包含文章id的）
	Body  string `json:"body" structs:"body"`   // 需要被搜索的正文
}

// Index
func (FullTextModel) Index() string {
	return "full_text_index"
}

func (FullTextModel) Mapping() string {
	return `
{
	"settings": {
		"index":{
			"max_result_window": "100000"
		}
	},
	"mappings": {
		"properties": {
			"key": {
				"type": "keyword"
			},
			"title": {
				"type": "text"
			},
			"slug": {
				"type": "keyword"
			},
			"body": {
				"type": "text"
			}
		}
	}
}
`
}

// IndexExists 判断索引是否存在
func (demo FullTextModel) IndexExists() bool {
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
func (demo FullTextModel) CreateIndex() error {
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
func (demo FullTextModel) RemoveIndex() error {
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
