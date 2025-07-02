package resp

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

func SendConnFailWithMsg(msg string, conn *websocket.Conn) {
	data := Response{FailValidCode, "错误", msg}
	byteData, _ := json.Marshal(data)
	conn.WriteMessage(websocket.TextMessage, byteData)
}

func SendConnOkWithData(data any, conn *websocket.Conn) {
	byteData, _ := json.Marshal(Response{
		Code: SuccessCode,
		Msg:  "成功",
		Data: data,
	})
	conn.WriteMessage(websocket.TextMessage, byteData)
}
func SendWsMsg(onLineMap map[uint]map[string]*websocket.Conn, userID uint, data any) {
	addrMap, ok := onLineMap[userID]
	if !ok {
		return
	}
	byteData, _ := json.Marshal(Response{SuccessCode, "成功", data})
	for _, conn := range addrMap {
		conn.WriteMessage(websocket.TextMessage, byteData)
	}
}
