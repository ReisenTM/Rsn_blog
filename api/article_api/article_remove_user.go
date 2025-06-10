package article_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

// ArticleRemoveUserView 用户不能批量删除
func (ArticleApi) ArticleRemoveUserView(c *gin.Context) {
	cr := middleware.GetBind[model.IDRequest](c)
	//拿到用户信息
	claim := jwts.GetClaims(c)
	//先找文章
	var article model.ArticleModel
	err := global.DB.Take(&article, "user_id = ? and id =?", claim.UserID, cr.ID).Error
	if err != nil {
		resp.FailWithMsg("文章不存在", c)
		return
	}
	err = global.DB.Delete(&article).Error
	if err != nil {
		resp.FailWithMsg(err.Error(), c)
		return
	}
	resp.OKWithMsg("删除成功", c)
}
