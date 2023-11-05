package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Params struct {
	ID   string `json:"id"`
	Size int    `json:"size"`
}

type NewData struct {
	Index    string `json:"index"`
	Title    string `json:"title"`
	HotValue string `json:"hotValue"`
	Link     string `json:"link"`
}

type NewResponse struct {
	Code int       `json:"code"`
	Data []NewData `json:"data"`
	Msg  string    `json:"msg"`
}

func main() {
	var params = Params{
		ID:   "mproPpoq6O",
		Size: 1,
	}
	reqParam, _ := json.Marshal(params)
	reqBody := strings.NewReader(string(reqParam))

	httpReq, err := http.NewRequest("POST", "https://api.codelife.cc/api/top/list", reqBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Signaturekey", "U2FsdGVkX19bXc4NrVJQzdIFTvvZ3UYN6uQlLBI4TCI=")
	httpReq.Header.Add("Version", "1.3.39")

	client := http.Client{
		Timeout: 2 * time.Second, // 两秒为期限
	}

	// DO: HTTP请求
	httpResp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println(err)
		return
	}

	var response NewResponse
	byteData, err := io.ReadAll(httpResp.Body)
	//fmt.Println(byteData)
	err = json.Unmarshal(byteData, &response)
	if err != nil {
		fmt.Println(err)
		return
	}

	if response.Code != 200 {
		fmt.Println("状态码非200", response.Msg)
		return
	}
	fmt.Println(response.Data)
}
