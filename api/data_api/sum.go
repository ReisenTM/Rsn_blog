package data_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/service/redis_service/redis_site"
	"github.com/gin-gonic/gin"
)

type SumResponse struct {
	FlowCount     int   `json:"flow_count"`
	UserCount     int64 `json:"user_count"`
	ArticleCount  int64 `json:"article_count"`
	MessageCount  int64 `json:"message-count"`
	CommentCount  int64 `json:"comment_count"`
	NewLoginCount int64 `json:"new_login_count"`
	NewSignCount  int64 `json:"new_sign_count"`
}

// SumView 数据和
func (DataApi) SumView(c *gin.Context) {
	var data SumResponse
	data.FlowCount = redis_site.GetFlow()
	global.DB.Model(model.UserModel{}).Count(&data.UserCount)
	global.DB.Model(model.ArticleModel{}).Where("status = ?", enum.ArticleStatusPublished).Count(&data.ArticleCount)
	global.DB.Model(model.ChatModel{}).Count(&data.MessageCount)
	global.DB.Model(model.CommentModel{}).Count(&data.CommentCount)
	global.DB.Model(model.UserLoginModel{}).Where("date(created_at) = date(now())").Count(&data.NewLoginCount)
	global.DB.Model(model.UserModel{}).Where("date(created_at) = date(now())").Count(&data.NewSignCount)
	resp.OkWithData(data, c)
}
