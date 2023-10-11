package global

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gvb_server/config"
)

var (
	Config   *config.Config // 这样才可以修改它
	DB       *gorm.DB       // 数据库
	Log      *logrus.Logger // 日志
	MysqlLog logger.Interface
	Redis    *redis.Client
)
