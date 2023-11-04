package es_ser

import "gvb_server/models"

type Option struct {
	models.PageInfo
	Fields []string
	// 后面如果key为空，那么就是
	// elastic.NewMultiMatchQuery(option.Key, "title", "abstract", "content")查询
	// 但是这个的返回值是[]string，所以重新定义一个Fields来接收这个
	Tag string
}

// GetForm 为了生效于原值，所以要用指针，即o *Option
func (o *Option) GetForm() int {
	if o.Page == 0 {
		o.Page = 1
	}
	if o.Limit == 0 {
		o.Limit = 10
	}
	return (o.Page - 1) * o.Limit
}
