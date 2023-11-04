package es_ser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/olivere/elastic/v7"
	"github.com/russross/blackfriday"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gvb_server/global"
	"gvb_server/models"
	"strings"
)

// SearchData 全文搜索的索引
type SearchData struct {
	Key   string `json:"key"`   // 文章关联的id，这样才能在根据文章title删除文章时删除全文搜素的文章
	Body  string `json:"body"`  // 需要被搜索的正文
	Slug  string `json:"slug"`  // 跳转地址（包含文章id的）
	Title string `json:"title"` // 标题
}

// AsyncArticleByFullText 批量同步到全局搜索的es中 -> 是异步操作 -> 不用返回err
func AsyncArticleByFullText(id, title, content string) {
	// 批量添加
	indexList := GetSearchIndexByContent(id, title, content)
	// 先创建一个桶
	bulk := global.ESClient.Bulk()
	for _, indexData := range indexList {
		req := elastic.NewBulkIndexRequest().Index(models.FullTextModel{}.Index()).Doc(indexData) // 即将加入桶的请求
		bulk.Add(req)
	}
	result, err := bulk.Do(context.Background()) // 执行，统一删除
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("%s 添加成功, 共 %d 条！", title, len(result.Succeeded()))
}

// 删除文章数据
func DeleteFullTextByArticleID(id string) {
	boolSearch := elastic.NewTermQuery("key", id)
	res, _ := global.ESClient.
		DeleteByQuery().
		Index(models.FullTextModel{}.Index()).
		Query(boolSearch).
		Do(context.Background())
	logrus.Infof("成功删除 %d 条记录", res.Deleted)
}

// GetSearchIndexByContent 根据文章内容获取文章对应的索引，索引格式为SearchData
func GetSearchIndexByContent(id, title string, content string) (searchDataList []SearchData) {
	// 按换行切割每一行
	dataList := strings.Split(content, "\n")
	// 区分是不是标题，即是否#开头
	// 但是代码块中的#不是，所以设bool值，true：在代码块中，false：不在代码块中，来区分
	var isCode bool = false
	var headList, bodyList []string // 标题和正文
	var body string
	headList = append(headList, getHeader(title)) // 标题

	for _, s := range dataList {
		// 同样可以用正则表达式来进行区分 -> #{1,6}
		// 判断是否是代码块
		if strings.HasPrefix(s, "```") {
			isCode = !isCode
		}
		if strings.HasPrefix(s, "#") && !isCode { // 是标题
			headList = append(headList, getHeader(s))
			bodyList = append(bodyList, getBody(body)) // 因为加了文章标题，所以这部分即使为空也可以放开
			body = ""
			continue
		}
		body += s
	}
	bodyList = append(bodyList, getBody(body))

	ln := len(headList)
	for i := 0; i < ln; i++ {
		searchDataList = append(searchDataList, SearchData{
			Title: headList[i],
			Body:  bodyList[i],
			Slug:  id + getSlug(headList[i]),
			Key:   id,
		})
	}

	//for _, se := range searchDataList {
	//	j, _ := json.Marshal(se)
	//	fmt.Println(string(j))
	//}
	return searchDataList
}

// getHeader 标题的格式优化
func getHeader(head string) string {
	// 把标题前面的"## "替换掉
	head = strings.ReplaceAll(head, "#", "")
	head = strings.ReplaceAll(head, " ", "")
	// 拼成前端格式"?id=" + head由前端自行弄
	return head
}

// getBody内容格式规范化
func getBody(body string) string {
	unsafe := blackfriday.MarkdownCommon([]byte(body))
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(unsafe)))
	return doc.Text()
}

// 获取跳转地址
func getSlug(slug string) string {
	return "#" + slug
}
