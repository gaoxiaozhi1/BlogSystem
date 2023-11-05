package new_api

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"gvb_server/models/res"
	"gvb_server/service/redis_ser"
	"gvb_server/utils/requests"
	"io"
	"time"
)

type params struct {
	ID   string `json:"id"`
	Size int    `json:"size"`
}

type header struct {
	SignatureKey string `form:"signaturekey" structs:"signaturekey"`
	Version      string `form:"version" structs:"version"`
	UserAgent    string `form:"User-Agent" structs:"User-Agent"`
}

type NewResponse struct {
	Code int                 `json:"code"`
	Data []redis_ser.NewData `json:"data"`
	Msg  string              `json:"msg"`
}

const newApi = "https://api.codelife.cc/api/top/list"
const timeout = 2 * time.Second

func (NewApi) NewListView(c *gin.Context) {
	var cr params
	var headers header
	err := c.ShouldBindJSON(&cr)
	err = c.ShouldBindHeader(&headers)
	if err != nil {
		res.FailWitheCode(res.ArgumentError, c)
		return
	}
	if cr.Size == 0 {
		cr.Size = 1
	}

	key := fmt.Sprintf("%s-%d", cr.ID, cr.Size)
	newsdata, _ := redis_ser.GetNews(key)
	if len(newsdata) > 0 {
		res.OKWithData(newsdata, c)
		return
	}
	httpResponse, err := requests.Post(newApi, cr, structs.Map(headers), timeout)
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}

	var response NewResponse
	byteData, err := io.ReadAll(httpResponse.Body)
	//fmt.Println(byteData)
	err = json.Unmarshal(byteData, &response)
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}

	if response.Code != 200 {
		res.FailWithMessage(response.Msg, c)
		return
	}
	res.OKWithData(response, c)
	redis_ser.SetNews(key, response.Data)
	return
}
