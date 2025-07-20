package resp

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// SSEOk 服务器-》客户端 单向通知
func SSEOk(data any, c *gin.Context) {
	byteData, _ := json.Marshal(Response{SuccessCode, data, "成功"})
	c.SSEvent("", string(byteData))
	c.Writer.Flush()
}

func SSEFail(msg string, c *gin.Context) {
	byteData, _ := json.Marshal(Response{FailServiceCodee, map[string]any{}, msg})
	c.SSEvent("", string(byteData))
	c.Writer.Flush()
}
