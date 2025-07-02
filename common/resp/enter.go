package resp

import (
	"blogX_server/utils/validate"
	"github.com/gin-gonic/gin"
	"net/http"
)

var empty = map[string]any{}

type Code int

func (c Code) String() string {
	switch c {
	case SuccessCode:
		return "成功"
	case FailValidCode:
		return "失败"
	}
	return ""
}

const (
	SuccessCode Code = iota
	FailValidCode
	FailServiceCodee
)

type Response struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (r *Response) Json(c *gin.Context) {
	c.JSON(http.StatusOK, r)
}

func OK(data any, msg string, c *gin.Context) {
	response := Response{Code: SuccessCode, Msg: msg, Data: data}
	response.Json(c)
}
func OKWithMsg(msg string, c *gin.Context) {
	response := Response{Code: SuccessCode, Msg: msg, Data: empty}
	response.Json(c)
}
func OkWithData(data any, c *gin.Context) {
	resp := Response{SuccessCode, "成功", data}
	resp.Json(c)
}
func OkWithList(list any, count int, c *gin.Context) {
	resp := Response{SuccessCode, "成功", map[string]any{
		"list":  list,
		"count": count,
	}}
	resp.Json(c)
}
func FailWithCode(code Code, c *gin.Context) {
	response := Response{Code: code, Msg: "失败", Data: empty}
	response.Json(c)
}
func FailWithMsg(msg string, c *gin.Context) {
	response := Response{Code: FailValidCode, Msg: msg, Data: empty}
	response.Json(c)
}
func FailWithData(data any, c *gin.Context) {
	response := Response{
		Code: SuccessCode,
		Msg:  "成功",
		Data: data,
	}
	response.Json(c)
}
func FailWithError(err error, c *gin.Context) {
	data, msg := validate.ValidateError(err)
	response := Response{Code: FailValidCode, Msg: msg, Data: data}
	response.Json(c)
}
