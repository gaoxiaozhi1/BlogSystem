package settings_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/config"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/models/res"
)

// 网站的相关信息，但是这个只能在配置文件修改，不太方便，后续最好优化成可以直接在前端页面修改的那种
// 配置文件的修改
func (SettingsApi) SettingsInfoUpdateView(c *gin.Context) {
	var cr config.SiteInfo
	err := c.ShouldBindJSON(&cr)

	if err != nil {
		res.FailWitheCode(res.ArgumentError, c) // 参数错误
		return
	}

	global.Config.SiteInfo = cr
	err = core.SetYaml()
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage(err.Error(), c)
		return
	}
	res.OKWith(c)
}
