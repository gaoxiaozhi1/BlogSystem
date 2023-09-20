package global

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gvb_server/config"
)

var (
	Config *config.Config // 这样才可以修改它
	DB     *gorm.DB
	Log    *logrus.Logger
)
