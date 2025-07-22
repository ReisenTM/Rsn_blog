package article_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"github.com/gin-gonic/gin"
)

type ArticleRecommendRequest struct {
	common.PageInfo
}

type ArticleRecommendResponse struct {
	ID         uint   `json:"id" gorm:"column:id"`
	Title      string `json:"title" gorm:"column:title"`
	ViewsCount int    `json:"views_count" gorm:"column:views_count"`
}

// ArticleRecommendView 今日热门文章推荐
func (ArticleApi) ArticleRecommendView(c *gin.Context) {
	cr := middleware.GetBind[ArticleRecommendRequest](c)
	if cr.Limit == 0 {
		cr.Limit = 6
	}
	var list = make([]ArticleRecommendResponse, 0)
	global.DB.Model(model.ArticleModel{}).Debug().
		Order("views_count desc").
		Where("date(created_at) = date(now())").
		Limit(cr.Limit).Select("id", "title", "views_count").Scan(&list)

	resp.OkWithList(list, len(list), c)
}
