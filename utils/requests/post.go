package requests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Post 通用的post请求
func Post(url string, data any, headers map[string]interface{}, timeout time.Duration) (resp *http.Response, err error) {
	reqParam, _ := json.Marshal(data)
	reqBody := strings.NewReader(string(reqParam))

	httpReq, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return
	}
	httpReq.Header.Add("Content-Type", "application/json") // 默认
	for k, v := range headers {
		// 因为必须要 传interface{}类型的
		switch val := v.(type) {
		case string:
			httpReq.Header.Add(k, val)
		case int:
			httpReq.Header.Add(k, strconv.Itoa(val))
		}
	}

	client := http.Client{
		Timeout: timeout, // 两秒为期限
	}
	// DO: HTTP请求
	httpResp, err := client.Do(httpReq)
	return httpResp, err
}
