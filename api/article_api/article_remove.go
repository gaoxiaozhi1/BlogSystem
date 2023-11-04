package article_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/es_ser"
)

type IDListRequest struct {
	IDList []string `json:"id_list"`
}

func (ArticleApi) ArticleRemoveView(c *gin.Context) {
	var cr IDListRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		global.Log.Error(err)
		res.FailWitheCode(res.ArgumentError, c)
		return
	}

	// ES文章删除
	// 这行代码创建了一个批量服务，用于在Elasticsearch中执行批量操作。
	// Index(models.ArticleModel{}.Index())指定了要在哪个索引上执行操作，
	// Refresh("true")表示在操作完成后刷新索引。
	bulkService := global.ESClient.Bulk().Index(models.ArticleModel{}.Index()).Refresh("true")

	// 如果文章删除了，用户收藏这篇文章怎么办?
	// 1．顺带把与这个文章关联的收藏也删除了
	// 2．用户收藏表，新增一个字段，表示文章是否删除，用户可以删除这个收藏记录，但是找不到文章去改收藏数
	for _, id := range cr.IDList {
		req := elastic.NewBulkDeleteRequest().Id(id) // 对于每个ID，创建一个新的批量删除请求。
		bulkService.Add(req)                         // 将删除请求添加到批量服务中。
		// 同步删除全局搜索es中的对应数据 -> 用 go 去删进度会快一点
		go es_ser.DeleteFullTextByArticleID(id)
	}
	// 执行批量服务中的所有请求，并返回结果和错误（如果有）
	result, err := bulkService.Do(context.Background())
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("删除失败", c)
		return
	}
	res.OKWithMessage(fmt.Sprintf("成功删除 %d 篇文章", len(result.Succeeded())), c)
	return
}
