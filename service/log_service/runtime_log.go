package log_service

import (
	"blogX_server/Model"
	"blogX_server/Model/enum"
	"blogX_server/global"
	"encoding/json"
	"fmt"
	e "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

// RuntimeLog 运行日志
type RuntimeLog struct {
	log             *Model.LogModel
	title           string
	itemList        []string //响应和请求中间的显示内容
	level           enum.LogLevel
	serviceName     string
	runtimeDateType RuntimeDateType
}

// Save 保存到数据库
func (r *RuntimeLog) Save() {
	r.SetNowTime()
	//判断创建还是更新
	var log Model.LogModel

	//找到满足条件的记录
	global.DB.Debug().Find(&log, fmt.Sprintf("service_name = ? and type = ? and created_at >= (now()-%s)",
		r.runtimeDateType.GetSqlTime()), r.serviceName, enum.LogRuntimeType)
	content := strings.Join(r.itemList, "\n")
	if log.ID != 0 {
		// 更新
		c := strings.Join(r.itemList, "\n")
		newContent := log.Content + "\n" + c

		// 之前已经save过了，那就是更新
		global.DB.Model(&log).Updates(map[string]any{
			"content": newContent,
		})
		r.itemList = []string{}
		return
	}
	err := global.DB.Create(&Model.LogModel{
		Type:        enum.LogRuntimeType,
		Title:       r.title,
		Content:     content,
		Level:       r.level,
		ServiceName: r.serviceName,
	}).Error
	if err != nil {
		logrus.Errorf("保存运行日志失败, %v", err)
		return
	}
}

func (r *RuntimeLog) SetTitle(title string) {
	r.title = title
}

func (r *RuntimeLog) SetLevel(level enum.LogLevel) {
	r.level = level
}
func (r *RuntimeLog) SetLink(label string, href string) {
	r.itemList = append(r.itemList, fmt.Sprintf("<div class=\"log_item link\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\"><a href=\"%s\" target=\"_blank\">%s</a></div></div> ",
		label,
		href, href))
}
func (r *RuntimeLog) SetImage(src string) {
	r.itemList = append(r.itemList, fmt.Sprintf("<div class=\"log_image\"><img src=\"%s\" alt=\"\"></div>", src))
}
func (r *RuntimeLog) setItem(label string, value any, logLevelType enum.LogLevel) {
	var v string

	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		byteData, _ := json.Marshal(value)
		v = string(byteData)
	default:
		v = fmt.Sprintf("%v", value)
	}

	r.itemList = append(r.itemList, fmt.Sprintf("<div class=\"log_item %s\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\">%s</div></div>",
		logLevelType,
		label, v))
}
func (r *RuntimeLog) SetItem(label string, value any) {
	r.setItem(label, value, enum.LogLevelInfo)
}
func (r *RuntimeLog) SetItemInfo(label string, value any) {
	r.setItem(label, value, enum.LogLevelInfo)
}
func (r *RuntimeLog) SetItemWarn(label string, value any) {
	r.setItem(label, value, enum.LogLevelWarn)
}
func (r *RuntimeLog) SetItemError(label string, value any) {
	r.setItem(label, value, enum.LogLevelError)
}

// SetNowTime 记录当前时间戳
func (r *RuntimeLog) SetNowTime() {
	r.itemList = append(r.itemList, fmt.Sprintf("<div class=\"log_time\">%s</div>", time.Now().Format("2006-01-02 15:04:05")))
}

func (r *RuntimeLog) SetError(label string, err error) {
	msg := e.WithStack(err)
	logrus.Errorf("%s %s", label, err.Error())
	r.itemList = append(r.itemList, fmt.Sprintf("<div class=\"log_error\"><div class=\"line\"><div class=\"label\">%s</div><div class=\"value\">%s</div><div class=\"type\">%T</div></div><div class=\"stack\">%+v</div></div>",
		label, err, err, msg))
}

// RuntimeDateType 运行日志时间：一小时内，一日内...
type RuntimeDateType uint8

const (
	RuntimeDateHour RuntimeDateType = iota + 1
	RuntimeDateDay
	RuntimeDateWeek
	RuntimeDateMonth
)

func (r RuntimeDateType) GetSqlTime() string {
	switch r {
	case RuntimeDateHour:
		return "interval '1 HOUR'"
	case RuntimeDateDay:
		return "interval '1 DAY'"
	case RuntimeDateWeek:
		return "interval '1 WEEK'"
	case RuntimeDateMonth:
		return "interval '1 MONTH'"
	}
	return "interval '1 DAY'"
}

// NewRuntimeLog 创建新运行日志
func NewRuntimeLog(serviceName string, dateType RuntimeDateType) *RuntimeLog {
	return &RuntimeLog{
		serviceName:     serviceName,
		runtimeDateType: dateType,
	}
}
