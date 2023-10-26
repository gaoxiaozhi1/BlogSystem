package main

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
	"strings"
)

func main() {
	// markdown 转 html
	unsafe := blackfriday.MarkdownCommon([]byte("### 你好\n ```python\nprint('你好')\n```\n - 123 \n \n<script>alert(123)</script>\n\n ![图片](http://xxx.com)"))
	fmt.Println(string(unsafe))
	// html 获取文本内容，xss过滤
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(unsafe)))
	//fmt.Println(doc.Text())
	// <script>alert(123)</script> 就属于恶意xss攻击
	// 如果要判断是否有xss攻击可以
	//nodes := doc.Find("script").Nodes
	//fmt.Println(nodes)

	doc.Find("script").Remove() // 把script标签删掉, xss过滤
	fmt.Println(doc.Text())

	// html 转 markdown
	converter := md.NewConverter("", true, nil)

	html, _ := doc.Html()
	markdown, err := converter.ConvertString(html)
	fmt.Println("md ->", markdown, err)
}
