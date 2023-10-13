package qq

import (
	"encoding/json"
	"fmt"
	"gvb_server/global"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// 用户信息
type QQInfo struct {
	Nickname string `json:"nickname"`     // 昵称
	Gender   string `json:"gender"`       // 性别
	Avatar   string `json:"figureurl_qq"` // 头像大图
	OpenID   string `json:"open_id"`
}

// QQ登录
type QQLogin struct {
	appID     string
	appKey    string
	redirect  string
	code      string
	accessTok string
	openID    string
}

// NewQQLogin 整个QQ登录的流程逻辑
func NewQQLogin(code string) (qqInfo QQInfo, err error) {
	qqLogin := &QQLogin{
		appID:    global.Config.QQ.AppID,
		appKey:   global.Config.QQ.Key,
		redirect: global.Config.QQ.Redirect,
		code:     code,
	}
	err = qqLogin.GetAccessToken()
	if err != nil {
		return qqInfo, err
	}
	err = qqLogin.GetOpenID()
	if err != nil {
		return qqInfo, err
	}
	qqInfo, err = qqLogin.GetUserInfo()
	if err != nil {
		return qqInfo, err
	}
	qqInfo.OpenID = qqLogin.openID
	return qqInfo, nil
}

// GetAccessToken 获取access_token
// 从QQ OAuth服务获取访问令牌的，这个访问令牌在以后可以用来访问和操作用户的账户数据。
// 并且QQLogin对象将包含新获取的访问令牌。
func (q *QQLogin) GetAccessToken() error {
	// 获取Access_token
	params := url.Values{} // 创建一个新的URL参数值
	// 这些行向URL参数添加了一些必要的字段
	// 如授权类型（grant_type）、客户端ID（client_id）、客户端密钥（client_secret）、授权码（code）和重定向URI（redirect_uri）。
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", q.appID)
	params.Add("client_secret", q.appKey)
	params.Add("code", q.code)
	params.Add("redirect_uri", q.redirect)

	// 解析URL字符串以创建一个新的URL对象。
	u, err := url.Parse("https://graph.qq.com/oauth2.0/token")
	if err != nil {
		return err
	}
	// 将之前创建和填充的参数编码为URL查询字符串，并将其设置为URL对象的原始查询部分。
	u.RawQuery = params.Encode()
	// 执行HTTP GET请求到指定的URL，并获取响应。
	res, err := http.Get(u.String())
	if err != nil {
		return err
	}

	defer res.Body.Close()       // // 延迟关闭响应体，直到包含此行代码的函数返回
	qs, err := parseQS(res.Body) // 解析响应体以获取查询字符串。
	if err != nil {
		return err
	}
	// 从查询字符串中获取访问令牌，并将其设置为QQLogin对象的访问令牌。
	q.accessTok = qs[`access_token`][0]
	return nil
}

// GetOpenID 获取openid
// 从QQ OAuth服务获取开放ID的，这个开放ID在以后可以用来访问和操作用户的账户数据。
// 这个开放ID是用户在QQ OAuth服务中的唯一标识符。
func (q *QQLogin) GetOpenID() error {
	// 获取openid
	// 解析URL字符串以创建一个新的URL对象，其中包含访问令牌（Access Token）。
	u, err := url.Parse(fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s", q.accessTok))
	if err != nil {
		return err
	}

	res, err := http.Get(u.String()) // 执行HTTP GET请求到指定的URL，并获取响应。
	if err != nil {
		return err
	}
	defer res.Body.Close() // 延迟关闭响应体，直到包含此行代码的函数返回。

	openID, err := getOpenID(res.Body) // 解析响应体以获取开放ID。
	if err != nil {
		return err
	}
	q.openID = openID // 将获取到的开放ID设置为QQLogin对象的开放ID。
	return nil
}

// GetUserInfo 获取用户信息
func (q *QQLogin) GetUserInfo() (qqInfo QQInfo, err error) {
	params := url.Values{}
	params.Add("access_token", q.accessTok)
	params.Add("oauth_consumer_key", q.appID)
	params.Add("openid", q.openID)
	u, err := url.Parse("https://graph.qq.com/user/get_user_info")
	if err != nil {
		return qqInfo, err
	}
	u.RawQuery = params.Encode()

	res, err := http.Get(u.String())
	data, err := io.ReadAll(res.Body)
	err = json.Unmarshal(data, &qqInfo)
	if err != nil {
		return qqInfo, err
	}
	return qqInfo, nil
}

// parseQS 将HTTP响应的正文解析为键值对的形式
func parseQS(r io.Reader) (val map[string][]string, err error) {
	val, err = url.ParseQuery(readAll(r))
	if err != nil {
		return val, err
	}
	return val, nil
}

// getOpenID 从HTTP响应的正文中解析出openid
func getOpenID(r io.Reader) (string, error) {
	body := readAll(r)
	start := strings.Index(body, `"openid":"`) + len(`"openid":"`)
	if start == -1 {
		return "", fmt.Errorf("openid not found")
	}
	end := strings.Index(body[start:], `"`)
	if end == -1 {
		return "", fmt.Errorf("openid not found")
	}
	return body[start : start+end], nil
}

// readAll 读取所有数据并将其转换为字符串
func readAll(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
