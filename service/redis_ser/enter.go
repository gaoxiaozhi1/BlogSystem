package redis_ser

const (
	articleDiggPrefix         = "digg_digg"
	articleLookPrefix         = "article_look"
	ArticleCommentCountPrefix = "article_comment_count"
	commentDiggPrefix         = "comment_digg"
)

func NewDigg() CountDB {
	return CountDB{
		Index: articleDiggPrefix,
	}
}
func NewArticleLook() CountDB {
	return CountDB{
		Index: articleLookPrefix,
	}
}
func NewCommentCount() CountDB {
	return CountDB{
		Index: ArticleCommentCountPrefix,
	}
}
func NewCommentDigg() CountDB {
	return CountDB{
		Index: commentDiggPrefix,
	}
}
