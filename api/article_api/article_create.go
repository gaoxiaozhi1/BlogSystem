package article_api

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"gvb_server/global"
	"gvb_server/models"
	"gvb_server/models/ctype"
	"gvb_server/models/res"
	"gvb_server/utils/jwts"
	"math/rand"
	"strings"
	"time"
)

type ArticleRequest struct {
	Title    string `json:"title" binding:"required" msg:"文章标题必填"`   // 文章标题
	Abstract string `json:"abstract"`                                // 文章简介
	Content  string `json:"content" binding:"required" msg:"文章内容必填"` // 文章内容
	Category string `json:"category"`                                // 文章分类
	Source   string `json:"source"`                                  // 文章来源
	Link     string `json:"link"`                                    // 原文连接
	// 后面便于不需要关联表查询，可以直接使用
	BannerID uint        `json:"banner_id"` // 文章封面id
	Tags     ctype.Array `json:"tags"`      // 文章标签
}

func (ArticleApi) ArticleCreateView(c *gin.Context) {
	var cr ArticleRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	userID := claims.UserID
	userNickName := claims.NickName
	// 校验content主要是防xss攻击

	// 处理content
	unsafe := blackfriday.MarkdownCommon([]byte(cr.Content))
	// html 获取文本内容，xss过滤
	// 是不是有script标签,即有没有xss恶意攻击
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(unsafe)))
	nodes := doc.Find("script").Nodes
	if len(nodes) > 0 {
		// 有script标签
		doc.Find("script").Remove() // 去掉标签
		// 重新转化成markdown格式
		converter := md.NewConverter("", true, nil)
		html, _ := doc.Html()
		markdown, _ := converter.ConvertString(html)
		cr.Content = markdown
	}

	// 截取文章简介
	if cr.Abstract == "" {
		abs := []rune(doc.Text())
		// 将content转为html，并且过滤xsS，以及获取正文内容
		if len(abs) > 100 {
			cr.Abstract = string(abs[0:100]) // 截取前100个字符
		} else {
			cr.Abstract = string(abs) // 截取所有
		}

	}

	// 如果不传banner_id 就从后台随机选一张
	if cr.BannerID == 0 {
		var bannerIDList []uint
		global.DB.Model(&models.BannerModel{}).Select("id").Scan(&bannerIDList)
		if len(bannerIDList) == 0 {
			res.FailWithMessage("没有banner数据", c)
			return
		}
		// 随机选择一个头像
		rand.Seed(time.Now().UnixNano()) // 初始化随机数生成器的种子的123
		cr.BannerID = bannerIDList[rand.Intn(len(bannerIDList))]
	}

	// 查banner_id下的banner_url
	var bannerUrl string
	err = global.DB.Model(&models.BannerModel{}).Where("id = ?", cr.BannerID).Select("path").Scan(&bannerUrl).Error
	if err != nil {
		res.FailWithMessage("banner不存在", c)
		return
	}

	// 查用户头像
	var avatar string
	err = global.DB.Model(&models.UserModel{}).Where("id = ?", userID).Select("avatar").Scan(&avatar).Error
	if err != nil {
		res.FailWithMessage("用户不存在", c)
		return
	}

	// 存储前的实例化
	now := time.Now().Format("2006-01-02 15:04:05")
	article := models.ArticleModel{
		CreatedAt:    now,
		UpdatedAt:    now,
		Title:        cr.Title,
		Abstract:     cr.Abstract,
		Content:      cr.Content,
		UserID:       userID,
		UserNickName: userNickName,
		UserAvatar:   avatar,
		Category:     cr.Category,
		Source:       cr.Source,
		Link:         cr.Link,
		BannerID:     cr.BannerID,
		BannerUrl:    bannerUrl,
		Tags:         cr.Tags,
	}

	err = article.Create()
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage(err.Error(), c)
		return
	}
	res.OKWithMessage("文章发布成功", c)
}
