package models

type CommentModel struct {
	//MODEL
	MODEL              `json:",select(c)"` // json-filter中tag为空 -> 匿名结构体为空
	SubComments        []CommentModel      `gorm:"foreignkey:ParentCommentID" json:"sub_comments,select(c)"` // 子评论列表
	ParentCommentModel *CommentModel       `gorm:"foreignkey:ParentCommentID" json:"comment_model"`          // 父级评论
	ParentCommentID    *uint               `json:"parent_comment_id,select(c)"`                              // 父级评论id
	Content            string              `gorm:"size:256" json:"content,select(c)"`                        // 评论内容
	DiggCount          int                 `gorm:"size:6;default:0" json:"digg_count,select(c)"`             // 点赞数
	CommentCount       int                 `gorm:"size:6;default:0" json:"comment_count,select(c)"`          // 子评论数
	ArticleID          string              `gorm:"size:32" json:"article_id,select(c)"`                      // 文章id
	User               UserModel           `json:"user,select(c)"`                                           // 关联的用户
	UserID             uint                `json:"user_id,select(c)"`                                        // 评论的用户id
}
