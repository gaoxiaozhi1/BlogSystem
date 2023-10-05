package res

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gvb_server/utils"
	"net/http"
)

// 封装返回的数据类型
type Response struct {
	Code int    `json:"code"` // 0：成功的；非0：失败的
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

type ListResponse[T any] struct {
	Count int64 `json:"count"`
	List  T     `json:"list"`
}

// 定义一个常量，(万一前端觉得0不是成功的)
const (
	Success = 0 // code为0：成功
	Error   = 7 // 其他错误
)

// 一些方法 ------- 返回给前端的是json数据
func Result(code int, data any, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// 封装一些相应成功的
func OK(data any, msg string, c *gin.Context) {
	Result(Success, data, msg, c)
}

func OKWithData(data any, c *gin.Context) {
	fmt.Println("这里的东西")
	Result(Success, data, "成功", c)
}

func OKWithList(list any, count int64, c *gin.Context) {
	OKWithData(ListResponse[any]{
		List:  list,
		Count: count,
	}, c)
}

func OKWithMessage(msg string, c *gin.Context) {
	Result(Success, map[string]any{}, msg, c)
}

func OKWith(c *gin.Context) {
	Result(Success, map[string]any{}, "成功", c)
}

// 封装一些相应失败的
func Fail(data any, msg string, c *gin.Context) {
	Result(Error, data, msg, c)
}
func FailWithMessage(msg string, c *gin.Context) {
	Result(Error, map[string]any{}, msg, c)
}

func FailWithError(err error, obj any, c *gin.Context) {
	// 获取报错中的信息
	msg := utils.GetValidMsg(err, obj)
	FailWithMessage(msg, c)
}

func FailWitheCode(code ErrorCode, c *gin.Context) {
	msg, ok := ErrorMap[code]
	if ok {
		Result(int(code), map[string]any{}, msg, c)
		return
	}
	Result(Error, map[string]any{}, "未知错误", c)

}
