package common

import (
	"gorm.io/gorm"
	"gvb_server/global"
	"gvb_server/models"
)

// 分页和搜索
type Option struct {
	models.PageInfo
	Debug bool // 如果为true，就看日志(即.Debug())
}

// 列表页
func ComList[T any](model T, option Option) (list []T, count int64, err error) {
	DB := global.DB
	if option.Debug {
		DB = global.DB.Session(&gorm.Session{Logger: global.MysqlLog})
	}

	// 查找图片数据库中的总数.RowsAffected
	count = DB.Select("id").Find(&list).RowsAffected
	// SELECT `id` FROM `banner_models`比
	// lobal.DB.Debug().Find(&imageList).RowsAffected的
	// SELECT * FROM `banner_models` 好

	offset := (option.Page - 1) * option.Limit // 偏移量，就是当前是从第几页开始的
	if offset < 0 {                            // 即cr.Page为0，那么offset为-1；
		offset = 0
	}
	// 分页查询（常用）
	err = DB.Limit(option.Limit).Offset(offset).Find(&list).Error

	return list, count, err
}
