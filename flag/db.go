package flag

import (
	"gvb_server/global"
	"gvb_server/models"
)

// 迁移表结构
// Makemigrations 初始化数据库表，数据库迁移等工作
func Makemigrations() {
	var err error
	// 这两行代码设置了两个联接表的结构。(多对多需要设置这个)
	// 设置用户和收藏之间的联接表的结构。这个联接表用于存储用户和他们的收藏之间的关系。
	// &models.UserModel{} 是你的用户模型，它代表了用户表。
	// "CollectsModels" 是你的收藏模型的名字，它代表了收藏表。
	// &models.UserCollectModel{} 是你的用户收藏模型，它代表了用户和收藏之间的联接表。
	//global.DB.SetupJoinTable(&models.UserModel{}, "CollectsModels", &models.UserCollectModel{})
	global.DB.SetupJoinTable(&models.MenuModel{}, "Banners", &models.MenuBannerModel{})

	// 这行代码使用 GORM 的 AutoMigrate 方法自动迁移表结构。
	// 这意味着，如果数据库中不存在这些表，那么 GORM 将会创建它们。
	// 如果这些表已经存在，但是它们的结构和你的 Go 结构体不一致，那么 GORM 将会更新这些表的结构以使其与 Go 结构体一致。
	// 这里需要注意的是，AutoMigrate 方法只会创建表、缺失的列和索引，它不会改变现有列的类型或删除未使用的列。
	// 生成四张表的表结构
	// 这行代码设置了 GORM 创建表时的选项。这里，它设置了 MySQL 的存储引擎为 InnoDB。
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(
			&models.BannerModel{},
			&models.TagModel{},
			&models.MessageModel{},
			&models.AdvertModel{},
			&models.UserModel{},
			&models.CommentModel{},
			//&models.ArticleModel{},
			&models.UserCollectModel{},
			&models.MenuModel{},
			&models.MenuBannerModel{},
			&models.FadeBackModel{},
			&models.LoginDataModel{},
			&models.ChatModel{},
		)

	if err != nil {
		global.Log.Error("[ error ] 生成数据库表结构失败")
	}
	global.Log.Infof("[ success ] 生成数据库表结构成功")
}
