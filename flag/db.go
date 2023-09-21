package flag

import (
	"gvb_server/global"
	"gvb_server/models"
)

// 迁移表结构
// Makemigrations 初始化数据库表，数据库迁移等工作
func Makemigrations() {
	var err error
	global.DB.SetupJoinTable(&models.UserModel{}, "CollectsModels", &models.UserCollectModel{})
	global.DB.SetupJoinTable(&models.MenuModel{}, "Banners", &models.MenuBannerModel{})
	// 生成四张表的表结构
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(
			&models.BannerModel{},
			&models.TagModel{},
			&models.MessageModel{},
			&models.AdvertModel{},
			&models.UserModel{},
			&models.CommentModel{},
			&models.ArticleModel{},
			&models.MenuModel{},
			&models.MenuBannerModel{},
			&models.FadeBackModel{},
			&models.LoginDataModel{},
		)

	if err != nil {
		global.Log.Error("[ error ] 生成数据库表结构失败")
	}
	global.Log.Infof("[ success ] 生成数据库表结构成功")
}
