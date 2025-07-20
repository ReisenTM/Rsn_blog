package comment_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/model/enum/relationship_enum"
	"blogX_server/service/focus_service"
	"blogX_server/service/redis_service/redis_comment"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
	"time"
)

type CommentListRequest struct {
	common.PageInfo
	ArticleID uint `form:"articleID"`
	UserID    uint `form:"userID"`
	Type      int8 `form:"type" binding:"required"` // 1 查我发文章的评论  2 查我发布的评论  3 管理员看所有的评论
}

type CommentListResponse struct {
	ID           uint                       `json:"id"`
	CreatedAt    time.Time                  `json:"createdAt"`
	Content      string                     `json:"content"`
	UserID       uint                       `json:"user_id"`
	UserNickname string                     `json:"user_nickname"`
	UserAvatar   string                     `json:"user_avatar"`
	ArticleID    uint                       `json:"article_id"`
	ArticleTitle string                     `json:"article_title"`
	ArticleCover string                     `json:"article_cover"`
	FavorCount   int                        `json:"favor_count"`
	Relation     relationship_enum.Relation `json:"relation,omitempty"`
	IsMe         bool                       `json:"is_me"`
}

// CommentListView 评论列表
func (CommentApi) CommentListView(c *gin.Context) {
	cr := middleware.GetBind[CommentListRequest](c)
	query := global.DB.Where("")
	claims := jwts.GetClaims(c)
	switch cr.Type {
	case 1:
		// 查我发文章的评论
		//查我发布的文章
		var articleIDList []uint
		global.DB.Model(model.ArticleModel{}).Where("user_id = ? and status = ?", claims.UserID, enum.ArticleStatusPublished).Select("id").Scan(&articleIDList)
		query.Where("article_id IN ?", articleIDList)
		cr.UserID = 0 //避免干扰
	case 2:
		// 查我发布的评论
		cr.UserID = claims.UserID
	case 3:
		//看所有评论
	}
	_list, count, _ := common.ListQuery(model.CommentModel{
		ArticleID: cr.ArticleID,
		UserID:    cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"UserModel", "ArticleModel"},
		Where:    query,
	})
	var RelationMao = map[uint]relationship_enum.Relation{}
	if cr.Type == 1 {
		var userIDList []uint
		for _, model := range _list {
			userIDList = append(userIDList, model.UserID)
		}
		RelationMao = focus_service.CalcUserPatchRelationship(claims.UserID, userIDList)
	}

	var list = make([]CommentListResponse, 0)
	for _, model := range _list {
		list = append(list, CommentListResponse{
			ID:           model.ID,
			CreatedAt:    model.CreatedAt,
			Content:      model.Content,
			UserID:       model.UserID,
			UserNickname: model.UserModel.Nickname,
			UserAvatar:   model.UserModel.Avatar,
			ArticleID:    model.ArticleID,
			ArticleTitle: model.ArticleModel.Title,
			ArticleCover: model.ArticleModel.Cover,
			FavorCount:   model.FavorCount + redis_comment.GetCacheFavor(model.ID),
			Relation:     RelationMao[model.UserID],
			IsMe:         model.UserID == claims.UserID,
		})
	}
	resp.OkWithList(list, count, c)

}
