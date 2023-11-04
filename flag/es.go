package flag

import "gvb_server/models"

// ES生成表结构

// EsCreateIndex es创建索引
func EsCreateIndex() {
	//models.ArticleModel{}.CreateIndex()
	models.FullTextModel{}.CreateIndex()
}

// EsRemoveIndex es删除索引
func EsRemoveIndex() {
	models.ArticleModel{}.RemoveIndex()
}
