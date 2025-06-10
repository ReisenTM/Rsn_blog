package article_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/service/redis_service/redis_article"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

type ArticleCollectRequest struct {
	ArticleID uint `json:"article_id" binding:"required"`
	CollectID uint `json:"collect_id"`
}

func (ArticleApi) ArticleCollectView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCollectRequest](c)

	var article model.ArticleModel
	err := global.DB.Take(&article, "status = ? and id = ?", enum.ArticleStatusPublished, cr.ArticleID).Error
	if err != nil {
		resp.FailWithMsg("文章不存在", c)
		return
	}
	var collectModel model.CollectModel
	claims := jwts.GetClaims(c)
	if cr.CollectID == 0 {
		// 是默认收藏夹
		err = global.DB.Take(&collectModel, "user_id = ? and is_default = ?", claims.UserID, 1).Error
		if err != nil {
			// 创建一个默认收藏夹
			collectModel.Title = "默认收藏夹"
			collectModel.UserID = claims.UserID
			collectModel.IsDefault = true
			global.DB.Create(&collectModel)
		}
		cr.CollectID = collectModel.ID
	} else {
		// 判断收藏夹是否存在，并且是否是自己创建的
		err = global.DB.Take(&collectModel, "user_id = ? ", claims.UserID).Error
		if err != nil {
			resp.FailWithMsg("收藏夹不存在", c)
			return
		}
	}

	// 判断是否收藏
	var articleCollect model.UserArticleCollectModel
	err = global.DB.Where(model.UserArticleCollectModel{
		UserID:    claims.UserID,
		ArticleID: cr.ArticleID,
		CollectID: cr.CollectID,
	}).Take(&articleCollect).Error

	if err != nil {
		// 收藏
		mod := model.UserArticleCollectModel{
			UserID:    claims.UserID,
			ArticleID: cr.ArticleID,
			CollectID: cr.CollectID,
		}
		err = global.DB.Create(&mod).Error
		if err != nil {
			resp.FailWithMsg("收藏失败", c)
			return
		}
		resp.OKWithMsg("收藏成功", c)

		redis_article.SetCacheCollect(cr.ArticleID, true)
		return
	}
	// 取消收藏
	err = global.DB.Where(model.UserArticleCollectModel{
		UserID:    claims.UserID,
		ArticleID: cr.ArticleID,
		CollectID: cr.CollectID,
	}).Delete(&articleCollect).Error
	if err != nil {
		resp.FailWithMsg("取消收藏失败", c)
		return
	}
	resp.OKWithMsg("取消收藏成功", c)
	redis_article.SetCacheCollect(cr.ArticleID, false)
	return
}
