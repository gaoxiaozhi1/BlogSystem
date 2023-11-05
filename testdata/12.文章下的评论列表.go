package main

import (
	"fmt"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/models"
)

func main() {
	// 读取配置文件
	core.InitConf()
	// 初始化日志
	global.Log = core.InitLogger()
	// 连接数据库
	global.DB = core.InitGorm()

	FindArticleCommentList("fPMzlYsBc9tzXF2Q8ZIH")
}

// FindArticleCommentList 根据文章id查找对应文章下的根评论
func FindArticleCommentList(articleID string) {
	// 先把文章下的根评论查出来
	var RootCommentList []*models.CommentModel
	global.DB.Find(&RootCommentList, "article_id = ? and parent_comment_id is null", articleID)
	//fmt.Println(RootCommentList)
	// 遍历根评论，递归查根评论下的所有子评论
	for _, model := range RootCommentList {
		//var subCommentList []models.CommentModel
		FindSubComment(*model, &model.SubComments)
		//model.SubComments = subCommentList
	}
	fmt.Println(RootCommentList[0])
}

// FindSubComment 递归查找评论下的子评论
func FindSubComment(model models.CommentModel, subCommentList *[]models.CommentModel) {
	global.DB.Preload("SubComments").Take(&model)
	//fmt.Println(*subCommentList)
	for _, sub := range model.SubComments {
		*subCommentList = append(*subCommentList, sub)
		FindSubComment(sub, subCommentList)
	}
	return
}
