package comment_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/service/comment_service"
	"blogX_server/service/message_service"
	"blogX_server/service/redis_service/redis_article"
	"blogX_server/service/redis_service/redis_comment"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

type CommentCreateRequest struct {
	Content   string `json:"content" binding:"required"`
	ArticleID uint   `json:"article_id" binding:"required"`
	ParentID  *uint  `json:"parent_id"` // 父评论id,存在没有的情况
}

// CommentCreateView 评论创建
func (CommentApi) CommentCreateView(c *gin.Context) {
	cr := middleware.GetBind[CommentCreateRequest](c)

	claim := jwts.GetClaims(c)

	var article model.ArticleModel
	err := global.DB.Take(&article, "user_id = ? and id = ?", claim.UserID, cr.ArticleID).Error
	if err != nil {
		resp.FailWithMsg("文章不存在", c)
		return
	}

	mod := model.CommentModel{
		Content:   cr.Content,
		UserID:    claim.UserID,
		ArticleID: cr.ArticleID,
		ParentID:  cr.ParentID,
	}
	//如果有父评论
	if cr.ParentID != nil {
		parentList := comment_service.GetParentComment(*cr.ParentID)
		//评论层级不能超过配置项
		if len(parentList) > global.Config.Site.Article.CommentLine {
			resp.FailWithMsg("评论层级达到限制", c)
			return
		}
		//找到根评论
		if len(parentList) > 0 {
			mod.RootID = &parentList[len(parentList)-1].ID
			for _, commentModel := range parentList {
				//回复数+1
				redis_comment.SetCacheReply(commentModel.ID, 1)
			}
			// 给父评论的用有人发消息
			defer func() {
				go message_service.InsertReplyMessage(mod)
			}()

		}
	}
	err = global.DB.Create(&mod).Error
	if err != nil {
		resp.FailWithMsg(err.Error(), c)
		return
	}
	//缓存评论数+1
	redis_article.SetCacheComment(cr.ArticleID, 1)
	resp.OKWithMsg("评论创建成功", c)
}
