package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// Mapping 映射
// "max_result_window": "100000" //列表最大能查到多少 --- 十万， 对于中型系统已经够啦
// "type": "text" // 可以进行模糊匹配。就是用于es搜索，可以高亮
// "type": "integer" // 可以比大小
// "format": "[yyyy-MM-dd HH:mm:ss]" // 时间，存储在es数据库里面的格式
func (DemoModel) Mapping() string {
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
		  "user_id": {
			"type": "integer" 
		  },
		  "created_at":{
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
func (demo DemoModel) IndexExists() bool {
	exists, err := client.
		IndexExists(demo.Index()).
		Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return exists
	}
	return exists
}

// CreateIndex 创建索引
func (demo DemoModel) CreateIndex() error {
	if demo.IndexExists() {
		// 有索引，就删掉索引，重新添加
		demo.RemoveIndex()
	}
	// 没有索引就创建索引
	createIndex, err := client.
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
func (demo DemoModel) RemoveIndex() error {
	logrus.Info("索引存在，删除索引")
	// 删除索引
	indexDelete, err := client.DeleteIndex(demo.Index()).Do(context.Background())
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
