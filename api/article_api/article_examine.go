package article_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
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

	//TODO:给发布者结果通知
	switch cr.Status {
	case enum.ArticleStatusPublished:

	case enum.ArticleStatusFail:

	}
	resp.OKWithMsg("文章审核通过", c)
}
