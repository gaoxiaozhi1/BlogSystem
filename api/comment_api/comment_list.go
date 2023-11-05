package comment_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/redis_ser"
)

type CommentListResponse struct {
	ArticleID string `form:"article_id"`
}

func (CommentApi) CommentListView(c *gin.Context) {
	var cr CommentListResponse
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	rootCommentList := FindArticleCommentList(cr.ArticleID)
	// json-filter空值问题，可以不显示隐私问题
	res.OKWithData(filter.Select("c", rootCommentList), c)
	return
}

// FindArticleCommentList 根据文章id查找对应文章下的根评论
func FindArticleCommentList(articleID string) (RootCommentList []*models.CommentModel) {
	// 先把文章下的根评论查出来
	global.DB.Preload("User").Find(&RootCommentList, "article_id = ? and parent_comment_id is null", articleID)
	//fmt.Println(RootCommentList)
	// 遍历根评论，递归查根评论下的所有子评论
	diggInfo := redis_ser.NewCommentDigg().GetInfo() // 存在redis缓存中的评论点赞数
	for _, model := range RootCommentList {
		var subCommentList, newSubCommentList []models.CommentModel
		FindSubComment(*model, &subCommentList)
		// 处理每一条评论的点赞数
		for _, commentModel := range subCommentList {
			digg := diggInfo[fmt.Sprintf("%d", commentModel.ID)]
			commentModel.DiggCount += digg
			newSubCommentList = append(newSubCommentList, commentModel)
		}
		// 处理根评论的点赞数
		modelDigg := diggInfo[fmt.Sprintf("%d", model.ID)]
		model.DiggCount += modelDigg
		model.SubComments = newSubCommentList
	}
	return
}

// FindSubComment 递归查找评论下的子评论
func FindSubComment(model models.CommentModel, subCommentList *[]models.CommentModel) {
	// Preload 预加载数据,可以把信息预加载且存储在对应model里面,拿的是子评论下的用户信息，所以
	global.DB.Preload("SubComments.User").Take(&model)
	//fmt.Println(*subCommentList)
	for _, sub := range model.SubComments {
		*subCommentList = append(*subCommentList, sub)
		FindSubComment(sub, subCommentList)
	}
	return
}

// FindSubCommentCount 递归查找评论下的子评论的数量
func FindSubCommentCount(model models.CommentModel) (subCommentList []models.CommentModel) {
	findSubCommentList(model, &subCommentList)
	return subCommentList
}

// findSubCommentList 递归查找子评论列表
func findSubCommentList(model models.CommentModel, subCommentList *[]models.CommentModel) {
	global.DB.Preload("SubComments").Take(&model)
	for _, sub := range model.SubComments {
		*subCommentList = append(*subCommentList, sub)
		FindSubComment(sub, subCommentList)
	}
	return
}
