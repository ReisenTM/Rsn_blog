package comment_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/service/message_service"
	"blogX_server/service/redis_service/redis_comment"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (CommentApi) CommentFavorView(c *gin.Context) {
	cr := middleware.GetBind[model.IDRequest](c)
	//找到评论
	var comment model.CommentModel
	err := global.DB.Take(&comment, cr.ID).Error
	if err != nil {
		resp.FailWithMsg("评论不存在", c)
		return
	}
	claims := jwts.GetClaims(c)
	var userFavorComment model.UserCommentFavorModel
	err = global.DB.Take(&userFavorComment, "user_id = ? AND comment_id = ?", claims.UserID, comment.ID).Error
	if err != nil {
		//如果之前没点赞过
		mod := model.UserCommentFavorModel{
			UserID:    claims.UserID,
			CommentID: cr.ID,
		}
		err = global.DB.Create(&mod).Error
		if err != nil {
			resp.FailWithMsg("点赞失败", c)
			logrus.Errorf("create comment favor")
			return
		}
		redis_comment.SetCacheFavor(cr.ID, 1)
		// 给这个评论的拥有人发消息
		message_service.InsertFavorCommentMessage(mod)

		resp.OKWithMsg("点赞成功", c)
		return
	}
	//如果已经点赞过
	err = global.DB.Delete(&userFavorComment).Error
	if err != nil {
		resp.FailWithMsg(err.Error(), c)
		logrus.Errorf("delete comment favor")
		return
	}
	redis_comment.SetCacheFavor(cr.ID, -1)
	resp.OKWithMsg("取消点赞成功", c)
}
