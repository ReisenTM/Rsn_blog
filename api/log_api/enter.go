package log_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/service/log_service"
	"fmt"
	"github.com/gin-gonic/gin"
)

type LogApi struct {
}

// LogListRequest 分页查询api结构体
type LogListRequest struct {
	common.PageInfo
	Type        enum.LogType  `form:"type"`
	UserID      uint          `json:"user_id"` //可能是访客，id设为0
	Level       enum.LogLevel `form:"level"`
	Title       string        `form:"title"`
	IP          string        `form:"ip"`
	LoginStatus bool          `form:"login_status"`
	ServiceName string        `form:"service_name"`
}
type LogListResponse struct {
	model.LogModel
	UserNickname string `json:"user_nickname"`
	UserAvatar   string `json:"user_avatar"`
}

// LogListView 分页查询
func (LogApi) LogListView(c *gin.Context) {
	var req LogListRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	list, count, err := common.ListQuery(model.LogModel{
		Type:        req.Type,
		Level:       req.Level,
		Title:       req.Title,
		UserID:      req.UserID,
		LoginStatus: false,
		ServiceName: req.ServiceName,
	}, common.Options{
		PageInfo:     req.PageInfo,
		Likes:        []string{"title"},
		Preloads:     []string{"UserModel"},
		Where:        nil,
		Debug:        false,
		DefaultOrder: "created_at desc",
	})
	//追加内容
	_list := make([]LogListResponse, 0)
	for _, item := range list {
		_list = append(_list, LogListResponse{
			LogModel:     item,
			UserNickname: item.UserModel.Nickname,
			UserAvatar:   item.UserModel.Avatar,
		})
	}
	resp.OkWithList(_list, int(count), c)
	return
}

// LogReadView 日志读取状态
func (LogApi) LogReadView(c *gin.Context) {
	var id model.IDRequest
	err := c.ShouldBindQuery(&id)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	var log model.LogModel
	err = global.DB.Take(&log, id).Error
	if err != nil {
		resp.FailWithMsg("不存在的日志", c)
		return
	}
	//如果找到了，设置为已读
	global.DB.Model(&log).Update("is_read", true)

	resp.OKWithMsg("日志读取成功", c)
}

// LogDeleteView 删除日志
func (LogApi) LogDeleteView(c *gin.Context) {
	var IDlist model.RemoveRequest
	err := c.ShouldBindJSON(&IDlist)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	//从中间件获取log对象,对删除操作写日志
	log := log_service.GetLog(c)
	log.SetShowReq()
	log.SetShowRes()

	var logs []model.LogModel
	err = global.DB.Where("id in (?)", IDlist.IDList).Find(&logs).Error
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	if len(logs) > 0 {
		//如果有数据
		global.DB.Delete(&logs)
	}
	msg := fmt.Sprintf("日志删除成功，共删除 %d 条日志", len(logs))
	resp.OKWithMsg(msg, c)
}
