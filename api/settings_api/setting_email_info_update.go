package settings_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/config"
	"gvb_server/core"
	"gvb_server/global"
	"gvb_server/models/res"
)

func (SettingsApi) SettingsEmailInfoUpdateView(c *gin.Context) {
	var cr config.Email
	err := c.ShouldBindJSON(&cr)

	if err != nil {
		res.FailWitheCode(res.ArgumentError, c) // 参数错误
		return
	}

	global.Config.Email = cr
	err = core.SetYaml()
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage(err.Error(), c)
		return
	}
	res.OKWith(c)
}
