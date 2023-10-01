package settings_api

import (
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models/res"
)

type SettingsUri struct {
	Name string `uri:"name"`
}

// 查看某一些配置信息（系统数据）
func (SettingsApi) SettingsInfoView(c *gin.Context) {

	var cr SettingsUri
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}
	switch cr.Name {
	case "site":
		res.OKWithData(global.Config.SiteInfo, c)
	case "email":
		res.OKWithData(global.Config.Email, c)
	case "qq":
		res.OKWithData(global.Config.QQ, c)
	case "qiniu":
		res.OKWithData(global.Config.QiNiu, c)
	case "jwt":
		res.OKWithData(global.Config.Jwy, c)
	default:
		res.FailWithMessage("没有对应的配置信息", c)
	}

}
