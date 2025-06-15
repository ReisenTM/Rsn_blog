package article_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/service/message_service"
	"blogX_server/service/redis_service/redis_article"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleFavorView(c *gin.Context) {
	cr := middleware.GetBind[model.IDRequest](c)
	var article model.ArticleModel
	err := global.DB.Take(&article, "id = ? and status = ?", cr.ID, enum.ArticleStatusPublished).Error
	if err != nil {
		resp.FailWithMsg("文章不存在", c)
		return
	}
	claims := jwts.GetClaims(c)
	// 查一下之前有没有点过
	var userFavorArticle model.UserArticleFavorModel
	err = global.DB.Take(&userFavorArticle, "user_id = ? and article_id = ?", claims.UserID, article.ID).Error
	if err != nil {
		// 点赞
		m := model.UserArticleFavorModel{
			UserID:    claims.UserID,
			ArticleID: cr.ID,
		}
		err = global.DB.Create(&m).Error
		if err != nil {
			resp.FailWithMsg("点赞失败", c)
			return
		}
		//给文章拥有者发消息
		message_service.InsertFavorArticleMessage(m)

		redis_article.SetCacheFavor(cr.ID, true)
		resp.OKWithMsg("点赞成功", c)
		return
	}
	// 点过就取消点赞
	global.DB.Delete(&userFavorArticle, "user_id = ? and article_id = ?", claims.UserID, article.ID)
	resp.OKWithMsg("取消点赞成功", c)
	redis_article.SetCacheFavor(cr.ID, false)
	return
}
