package article_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/ctype"
	"blogX_server/model/enum"
	"blogX_server/utils"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

// ArticleTagOptionsView 文章标签选择列表接口，方便前端调用
func (ArticleApi) ArticleTagOptionsView(c *gin.Context) {
	claims := jwts.GetClaims(c)

	var articleList []model.ArticleModel
	global.DB.Find(&articleList, "user_id = ? and status = ?", claims.UserID, enum.ArticleStatusPublished)

	var tagList ctype.List
	for _, m := range articleList {
		tagList = append(tagList, m.Tags...)
	}
	tagList = utils.Unique(tagList)
	var list = make([]model.OptionsResponse[string], 0)
	for _, s := range tagList {
		list = append(list, model.OptionsResponse[string]{
			Label: s,
			Value: s,
		})
	}
	resp.OkWithData(list, c)
}
