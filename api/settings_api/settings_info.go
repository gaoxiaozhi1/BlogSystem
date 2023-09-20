package settings_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/models/res"
)

// 系统数据
func (SettingsApi) SettingsInfoView(c *gin.Context) {
	res.OK(map[string]string{}, "xxx", c)
}
