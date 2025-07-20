package chat_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/ctype/chat_msg"
	"blogX_server/model/enum/chat_msg_type"
	"blogX_server/model/enum/relationship_enum"
	"blogX_server/service/focus_service"
	"blogX_server/utils/jwts"
	"blogX_server/utils/xss"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var UP = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ChatRequest struct {
	RevUserID uint                  `json:"rev_user_id"` // 发给谁
	MsgType   chat_msg_type.MsgType `json:"msg_type"`    // 1 文本 2 图片  3 md
	Msg       chat_msg.ChatMsg      `json:"msg"`         // 消息主体
}
type ChatResponse struct {
	ChatListResponse
}

// OnlineMap 在线map
var OnlineMap = map[uint]map[string]*websocket.Conn{} //用户->聊天关系映射

func (ChatApi) ChatView(c *gin.Context) {
	res := c.Writer
	req := c.Request
	// 服务升级
	conn, err := UP.Upgrade(res, req, nil)
	if err != nil {
		logrus.Errorf("websocket upgrade error: %s", err.Error())
		return
	}
	//解析token
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		logrus.Errorf("jwt parse error: %s", err.Error())
		return
	}
	userID := claims.UserID
	var user model.UserModel
	err = global.DB.Find(&user, userID).Error
	if err != nil {
		logrus.Errorf("find user error: %s", err.Error())
		return
	}
	addr := conn.RemoteAddr().String()
	addrMap, ok := OnlineMap[userID]
	if !ok {
		//上线
		OnlineMap[userID] = map[string]*websocket.Conn{
			addr: conn,
		}
	} else {
		_, ok1 := addrMap[addr]
		if !ok1 {
			OnlineMap[userID][addr] = conn
		}
	}

	for {
		// 消息类型，消息，错误
		_, p, err := conn.ReadMessage()
		if err != nil {
			// 一般是客户端断开 // websocket: close 1005 (no status)
			break
		}
		//json解码
		var req ChatRequest
		err = json.Unmarshal(p, &req)
		if err != nil {
			fmt.Println("收到原始消息:", string(p))
			logrus.Errorf("json unmarshal error: %s", err.Error())
			resp.SendConnFailWithMsg("请求解码失败", conn)
			return
		}
		// 判断接收人在不在
		var revUser model.UserModel
		err = global.DB.Take(&revUser, req.RevUserID).Error
		if err != nil {
			resp.SendConnFailWithMsg("接收人不存在", conn)
			continue
		}
		//根据消息类型做出决策
		switch req.MsgType {
		case chat_msg_type.TextMsgType:
			if req.Msg.TextMsg == nil || req.Msg.TextMsg.Content == "" {
				resp.SendConnFailWithMsg("文本消息内容为空", conn)
				continue
			}
		case chat_msg_type.ImageMsgType:
			if req.Msg.ImageMsg == nil || req.Msg.ImageMsg.Src == "" {
				resp.SendConnFailWithMsg("图片消息内容为空", conn)
				continue
			}
		case chat_msg_type.MarkdownMsgType:
			if req.Msg.MarkdownMsg == nil || req.Msg.MarkdownMsg.Content == "" {
				resp.SendConnFailWithMsg("markdown消息内容为空", conn)
				continue
			}
			// 对markdown消息做过滤
			req.Msg.MarkdownMsg.Content = xss.XssFilter(req.Msg.MarkdownMsg.Content)
		default:
			resp.SendConnFailWithMsg("不支持的消息类型", conn)
			continue
		}
		//判断双方关系，如果没有互相关注只能发三条消息
		// 陌生人，如果对方开了陌生人私信，那么就能聊
		//如果
		relation := focus_service.CalcUserRelationship(userID, req.RevUserID)
		switch relation {
		case relationship_enum.RelationStranger:
			//如果是陌生人
			var revUserMsgConf model.UserMessageConfModel
			err = global.DB.Take(&revUserMsgConf, "user_id = ?", revUser.ID).Error
			if err != nil {
				resp.SendConnFailWithMsg("接收人隐私设置不存在", conn)
				continue
			}
			if !revUserMsgConf.OpenPrivateChat {
				resp.SendConnFailWithMsg("对方未开始陌生人私聊", conn)
				continue
			}
		case relationship_enum.RelationFocus, relationship_enum.RelationFans: // 已关注
			// 今天对方如果没有回复你，那么你就只能发一条
			var chatList []model.ChatModel
			global.DB.Find(&chatList, "date(created_at) = date (now()) and ( (send_user_id = ? and  rev_user_id = ?) or (send_user_id = ? and  rev_user_id = ?))",
				userID, req.RevUserID, req.RevUserID, userID)

			// 我发的  对方发的
			var sendChatCount, revChatCount int
			for _, chat := range chatList {
				if chat.SendUserID == userID {
					sendChatCount++
				}
				if chat.RevUserID == userID {
					revChatCount++
				}
			}
			fmt.Println(sendChatCount, revChatCount)
			if sendChatCount > 1 && revChatCount == 0 {
				resp.SendConnFailWithMsg("对方未回复的情况下，当天只能发送一条消息", conn)
				continue
			}

		}
		//落库
		var chatModel model.ChatModel = model.ChatModel{
			SendUserID: userID,
			RevUserID:  req.RevUserID,
			MsgType:    req.MsgType,
			Msg:        req.Msg,
		}
		err = global.DB.Create(&chatModel).Error
		if err != nil {
			resp.SendConnFailWithMsg("消息发送失败", conn)
			continue
		}
		//填充返回体
		item := ChatResponse{ChatListResponse{
			ChatModel:        chatModel,
			SendUserNickname: user.Nickname,
			SendUserAvatar:   user.Avatar,
			RevUserNickname:  revUser.Nickname,
			RevUserAvatar:    revUser.Avatar,
		}}
		//给对方发一份，给自己发一份
		resp.SendWsMsg(OnlineMap, req.RevUserID, item)
		// 发给自己的,打上isme标签用于区分
		item.IsMe = true
		resp.SendConnOkWithData(item, conn)

	}
	defer func() {
		_ = conn.Close()
	}()
	addrMap2, ok2 := OnlineMap[userID]
	if ok2 {
		_, ok3 := addrMap2[addr]
		if ok3 {
			delete(OnlineMap[userID], addr)
		}
		if len(addrMap2) == 0 {
			delete(OnlineMap, userID)
		}
	}
}
