package core

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gvb_server/global"
	"time"
)

func InitGorm() *gorm.DB {
	if global.Config.Mysql.Host == "" {
		global.Log.Warnln("未配置mysql，取消gorm连接")
		return nil
	}
	dsn := global.Config.Mysql.Dsn()

	var mysqlLogger logger.Interface
	// 判断环境
	if global.Config.System.Env == "debug" {
		// 开发环境显示所有的sql
		mysqlLogger = logger.Default.LogMode(logger.Info) // 打印所有日志
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error) // 打印错误日志
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger的配置
		Logger: mysqlLogger,
	})

	// 连接失败就可以退出啦
	if err != nil {
		global.Log.Fatal(fmt.Sprintf("[%s] mysql连接失败"), dsn)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)               // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)              // 最多可容纳
	sqlDB.SetConnMaxLifetime(time.Hour * 4) // 连接最大复用时间，不能超过mysql的wait_timeout
	return db
}
