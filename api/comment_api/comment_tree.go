package comment_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/service/comment_service"
	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentTreeView(c *gin.Context) {
	cr := middleware.GetBind[model.IDRequest](c)

	var article model.ArticleModel
	//只允许已发布文章
	err := global.DB.Take(&article, "status = ? and id = ?", enum.ArticleStatusPublished, cr.ID).Error
	if err != nil {
		resp.FailWithMsg("文章不存在", c)
		return
	}
	// 把根评论查出来
	var commentList []model.CommentModel
	global.DB.Order("created_at desc").Find(&commentList, "article_id = ? and parent_id is null", cr.ID)
	var list = make([]comment_service.CommentResponse, 0)
	for _, mod := range commentList {
		response := comment_service.GetCommentTreeV4(mod.ID)
		list = append(list, *response)
	}
	resp.OkWithList(list, len(list), c)
}
