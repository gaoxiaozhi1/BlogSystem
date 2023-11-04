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
	// 是否看日志
	if option.Debug {
		DB = global.DB.Session(&gorm.Session{Logger: global.MysqlLog})
	}

	// 排序设置
	if option.Sort == "" {
		option.Sort = "created_at desc" // 默认按照时间往前排（从后往前排）
		//option.Sort = "created_at asc" // 默认按照时间往后排
	}

	//query := DB.Where(model)
	// 查找图片数据库中的总数.RowsAffected
	//count = DB.Where(model).Find(&list).RowsAffected
	//DB.Model(model).Where(model).Count(&count)
	//query = DB.Where(model)                    // 这里的query会受上面查询的影响，需要手动复位
	offset := (option.Page - 1) * option.Limit // 偏移量，就是当前是从第几页开始的
	if offset < 0 {                            // 即cr.Page为0，那么offset为-1；
		offset = 0
	}
	// 分页查询（常用）
	err = DB.Model(model).Where(model).
		Limit(option.Limit).Offset(offset).Order(option.Sort).
		Find(&list).Count(&count).Error

	return list, count, err
}
