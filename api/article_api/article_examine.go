package article_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/service/message_service"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ArticleExamineRequest struct {
	ArticleID uint               `json:"article_id" binding:"required"`
	Status    enum.ArticleStatus `json:"status" binding:"required,oneof=3 4"`
	Msg       string             `json:"msg"` //反馈信息
}

func (ArticleApi) ArticleExamineView(c *gin.Context) {
	cr := middleware.GetBind[ArticleExamineRequest](c)
	var article model.ArticleModel
	err := global.DB.Take(&article, cr.ArticleID).Error
	if err != nil {
		resp.FailWithMsg("文章不存在", c)
		return
	}
	//更新审核状态
	global.DB.Model(&article).Update("status", cr.Status)

	switch cr.Status {
	case enum.ArticleStatusPublished:
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", "审核成功", article.Title, fmt.Sprintf("/article/%d", article.ID))

	case enum.ArticleStatusFail:
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", fmt.Sprintf("审核失败 失败原因：%s", cr.Msg), "", "")
	}
	resp.OKWithMsg("文章审核通过", c)
}
