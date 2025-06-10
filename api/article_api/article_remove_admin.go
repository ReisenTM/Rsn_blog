package article_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

// ArticleRemoveAdminView 管理员可以批量删除
func (ArticleApi) ArticleRemoveAdminView(c *gin.Context) {
	cr := middleware.GetBind[model.RemoveRequest](c)
	var articles []model.ArticleModel
	err := global.DB.Find(&articles, "id in ?", cr.IDList).Error
	if err != nil {
		resp.FailWithMsg("文章不存在", c)
		return
	}
	if len(articles) > 0 {
		//for _, model := range articles {
		//TODO:给被删除的发通知[通知]
		//}
		err := global.DB.Delete(&articles).Error
		if err != nil {
			resp.FailWithMsg("删除失败", c)
			return
		}
	}

	resp.OKWithMsg(fmt.Sprintf("删除成功 成功删除%d条", len(articles)), c)
}
