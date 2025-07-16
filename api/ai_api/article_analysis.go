package ai_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/service/ai_service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
)

type ArticleAnalysisRequest struct {
	Content string `json:"content" binding:"required"`
}

type ArticleAnalysisResponse struct {
	Title    string   `json:"title"`
	Abstract string   `json:"abstract"`
	Category string   `json:"category"`
	Tag      []string `json:"tag"`
}

func (AIApi) ArticleAnalysisView(c *gin.Context) {
	cr := middleware.GetBind[ArticleAnalysisRequest](c)

	if !global.Config.AI.Enable {
		resp.FailWithMsg("站点未启用ai功能", c)
		return
	}

	md, err := ai_service.DSToChat(cr.Content)
	if err != nil {
		logrus.Errorf("ai分析失败 %s %s", err, cr.Content)
		resp.FailWithMsg("ai分析失败", c)
		return
	}
	strSlice := strings.Split(md, "```")
	msg := strSlice[1]
	after, _ := strings.CutPrefix(msg, "json")

	var data ArticleAnalysisResponse
	err = json.Unmarshal([]byte(after), &data)
	if err != nil {
		logrus.Errorf("ai分析失败 %s %s", err, msg)
		resp.FailWithMsg("ai分析失败", c)
		return
	}

	resp.OkWithData(data, c)
}
