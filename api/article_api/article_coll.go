package article_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/res"
	"gvb_server/service/es_ser"
	"gvb_server/utils/jwts"
)

// ArticleCollCreateView 用户收藏文章或取消收藏文章
func (ArticleApi) ArticleCollCreateView(c *gin.Context) {
	var cr models.ESIDRequest // 根据文章id查
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	model, err := es_ser.CommDetail(cr.ID) // 查文章详情
	if err != nil {
		res.FailWithMessage("文章不存在", c)
		return
	}

	var coll models.UserCollectModel
	err = global.DB.Take(&coll, "user_id = ? and article_id = ?", claims.UserID, cr.ID).Error
	//fmt.Println(coll.UserID, coll.ArticleID)
	fmt.Println(err)
	var num = -1
	if err != nil {
		// 没有找到 就收藏文章
		global.DB.Create(&models.UserCollectModel{
			UserID:    claims.UserID,
			ArticleID: cr.ID,
		})
		// 给文章的收藏数 +1
		num = 1
	}
	//fmt.Println(num)
	// 取消收藏
	// 给文章的收藏数 -1
	global.DB.Delete(&coll)

	// 更新文章收藏数
	err = es_ser.ArticleUpdate(cr.ID, map[string]any{
		"collects_count": model.CollectsCount + num,
	})

	if num == 1 {
		res.OKWithMessage("收藏文章成功", c)
	} else {
		res.OKWithMessage("取消收藏成功", c)
	}

}
