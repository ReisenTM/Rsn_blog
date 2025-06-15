package comment_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/service/comment_service"
	"blogX_server/service/message_service"
	"blogX_server/service/redis_service/redis_comment"
	"blogX_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (CommentApi) CommentRemoveView(c *gin.Context) {
	cr := middleware.GetBind[model.IDRequest](c)
	claim := jwts.GetClaims(c)
	var comment model.CommentModel
	err := global.DB.Preload("ArticleModel").Take(&comment, cr.ID).Error
	if err != nil {
		resp.FailWithMsg("评论不存在", c)
		return
	}
	if claim.Role != enum.RoleAdminType {
		// 普通用户只能删自己发的评论，或者自己发的文章的评论
		if !(comment.UserID == claim.UserID || comment.ArticleModel.UserID == claim.UserID) {
			resp.FailWithMsg("权限错误", c)
			return
		}
	}
	//通知用户评论被删除
	message_service.InsertSystemMessage(comment.UserID, "管理员删除了你的评论", fmt.Sprintf("%s 内容不符合社区规范", comment.Content), "", "")
	// 找所有的子评论，还要找所有的父评论
	subList := comment_service.GetCommentOneDimensional(comment.ID)
	if comment.ParentID != nil {
		// 有父评论
		parentList := comment_service.GetParentComment(*comment.ParentID)
		for _, commentModel := range parentList {
			redis_comment.SetCacheReply(commentModel.ID, -len(subList))
		}
	}
	err = global.DB.Delete(&subList).Error
	if err != nil {
		resp.FailWithError(err, c)
		logrus.Errorf(err.Error())
		return
	}
	msg := fmt.Sprintf("删除成功，共删除评论%d条", len(subList))
	resp.OKWithMsg(msg, c)
}
