package comment_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/redis_ser"
	"gvb_server/utils"
)

func (CommentApi) CommentRemoveView(c *gin.Context) {
	var cr CommentIDRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}
	var commentModel models.CommentModel
	err = global.DB.Take(&commentModel, cr.ID).Error
	if err != nil {
		res.FailWithMessage("评论不存在", c)
		return
	}
	// 统计评论下的子评论数 再把自己算上
	subCommentList := FindSubCommentCount(commentModel)
	count := len(subCommentList) + 1
	redis_ser.NewCommentCount().SetCount(commentModel.ArticleID, -count)

	// 判断是否是子评论
	if commentModel.ParentCommentID != nil {
		// 子评论 -> 那就要找其父评论, 减掉对应的评论数
		global.DB.Model(&models.CommentModel{}).
			Where("id = ?", *commentModel.ParentCommentID).
			Update("comment_count", gorm.Expr("comment_count - ?", count))
	}

	// 删除子评论以及当前评论
	var deleteCommentIDList []uint
	for _, model := range subCommentList {
		deleteCommentIDList = append(deleteCommentIDList, model.ID)
	}
	// 反转之后一个一个删除，避免comment外键检查
	utils.Reverse(deleteCommentIDList)
	deleteCommentIDList = append(deleteCommentIDList, commentModel.ID) // 删除当前评论
	for _, id := range deleteCommentIDList {                           // 一个个删除
		global.DB.Model(models.CommentModel{}).Delete("id = ?", id)
	}

	res.OKWithMessage(fmt.Sprintf("共删除了 %d 条评论", count), c)

}
