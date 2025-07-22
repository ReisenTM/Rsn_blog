package comment_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/model/enum/relationship_enum"
	"blogX_server/service/comment_service"
	"blogX_server/service/focus_service"
	"blogX_server/utils"
	"blogX_server/utils/jwts"
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

	var userRelationMap = map[uint]relationship_enum.Relation{}
	var userFavorCommentMap = map[uint]bool{}
	claims, err := jwts.ParseTokenByGin(c)
	if err == nil && claims != nil {
		// 登录了
		var commentList []model.CommentModel // 文章的评论id列表
		global.DB.Find(&commentList, "article_id = ?", cr.ID)

		if len(commentList) > 0 {
			// 查我点赞的评论id列表
			var commentIDList []uint
			var userIDList []uint
			for _, model := range commentList {
				commentIDList = append(commentIDList, model.ID)
				userIDList = append(userIDList, model.UserID)
			}
			userIDList = utils.Unique(userIDList) // 对用户id去重
			userRelationMap = focus_service.CalcUserPatchRelationship(claims.UserID, userIDList)
			var commentFavorList []model.UserCommentFavorModel
			global.DB.Find(&commentFavorList, "user_id = ? and comment_id in ?", claims.UserID, commentIDList)
			for _, model := range commentFavorList {
				userFavorCommentMap[model.CommentID] = true
			}
		}
	}
	// 把根评论查出来
	var commentList []model.CommentModel
	global.DB.Order("created_at desc").Find(&commentList, "article_id = ? and parent_id is null", cr.ID)
	var list = make([]comment_service.CommentResponse, 0)
	for _, mod := range commentList {
		response := comment_service.GetCommentTreeV4(mod.ID, userRelationMap, userFavorCommentMap)
		list = append(list, *response)
	}
	resp.OkWithList(list, len(list), c)
}
