package comment_service

import (
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/service/redis_service/redis_comment"
	"time"
)

// GetRootComment GetRoots 递归找根评论
func GetRootComment(commentID uint) *model.CommentModel {
	//先找到评论
	var comment model.CommentModel
	err := global.DB.Take(&comment, commentID).Error
	if err != nil {
		return nil
	}
	if comment.ParentID == nil {
		//如果已经是顶级评论
		return &comment
	}
	//不是就递归找
	return GetRootComment(*comment.ParentID)
}

// GetCommentTree 获取评论树(原位修改)
func GetCommentTree(mdl *model.CommentModel) {
	err := global.DB.Preload("SubComment").Take(mdl).Error
	if err != nil {
		return
	}
	for _, subComment := range mdl.SubComment {
		GetCommentTree(subComment)
	}
}

// GetCommentTreeV2 获取评论树
func GetCommentTreeV2(id uint) (mod *model.CommentModel) {
	mod = &model.CommentModel{
		Model: model.Model{ID: id},
	}

	global.DB.Preload("SubComment").Take(mod)
	for i := 0; i < len(mod.SubComment); i++ {
		commentModel := mod.SubComment[i]
		item := GetCommentTreeV2(commentModel.ID)
		mod.SubComment[i] = item
	}
	return
}

type CommentResponse struct {
	ID         uint               `json:"id"`
	CreatedAt  time.Time          `json:"createdAt"`
	Nickname   string             `json:"nickname"`
	Avatar     string             `json:"avatar"`
	Content    string             `json:"content"`
	UserID     uint               `json:"user_id"`
	ArticleID  uint               `json:"article_id"`
	ParentID   *uint              `json:"parent_id"`
	SubComment []*CommentResponse `json:"sub_comment"`
	FavorCount int                `json:"favor_count"` //点赞数
	ReplyCount int                `json:"reply_count"` //回复数
}

func GetCommentTreeV3(id uint) (res *CommentResponse) {
	mod := &model.CommentModel{
		Model: model.Model{ID: id},
	}

	global.DB.Preload("UserModel").Preload("SubComment").Take(mod)
	res = &CommentResponse{
		ID:         mod.ID,
		CreatedAt:  mod.CreatedAt,
		Content:    mod.Content,
		UserID:     mod.UserID,
		Nickname:   mod.UserModel.Nickname,
		Avatar:     mod.UserModel.Avatar,
		ArticleID:  mod.ArticleID,
		ParentID:   mod.ParentID,
		FavorCount: mod.FavorCount + redis_comment.GetCacheFavor(mod.ID),
		ReplyCount: redis_comment.GetCacheReply(mod.ID),
		SubComment: make([]*CommentResponse, 0),
	}
	for _, commentModel := range mod.SubComment {
		res.SubComment = append(res.SubComment, GetCommentTreeV3(commentModel.ID))
	}
	return
}

func GetCommentTreeV4(id uint) (res *CommentResponse) {
	return getCommentTreeV4(id, 1)
}

//line 层级
func getCommentTreeV4(id uint, line int) (res *CommentResponse) {
	mod := &model.CommentModel{
		Model: model.Model{ID: id},
	}

	global.DB.Preload("UserModel").Preload("SubComment").Take(mod)
	res = &CommentResponse{
		ID:         mod.ID,
		CreatedAt:  mod.CreatedAt,
		Content:    mod.Content,
		UserID:     mod.UserID,
		Nickname:   mod.UserModel.Nickname,
		Avatar:     mod.UserModel.Avatar,
		ArticleID:  mod.ArticleID,
		ParentID:   mod.ParentID,
		FavorCount: mod.FavorCount + redis_comment.GetCacheFavor(mod.ID),
		ReplyCount: mod.ReplyCount + redis_comment.GetCacheReply(mod.ID),
		SubComment: make([]*CommentResponse, 0),
	}
	if line >= global.Config.Site.Article.CommentLine {
		return
	}
	for _, commentModel := range mod.SubComment {
		res.SubComment = append(res.SubComment, getCommentTreeV4(commentModel.ID, line+1)) //递归一次层级+1
	}
	return
}

// GetParentComment 获取所有父节点
func GetParentComment(commentID uint) (list []model.CommentModel) {
	var comment model.CommentModel
	err := global.DB.Take(&comment, commentID).Error
	if err != nil {
		return
	}
	list = append(list, comment)
	if comment.ParentID != nil {
		list = append(list, GetParentComment(*comment.ParentID)...)
	}
	return
}

// GetCommentOneDimensional 评论树一维化,用来找子评论
func GetCommentOneDimensional(id uint) (list []model.CommentModel) {
	mod := model.CommentModel{
		Model: model.Model{ID: id},
	}

	global.DB.Preload("SubComment").Take(&mod)
	list = append(list, mod)
	for _, commentModel := range mod.SubComment {
		_list := GetCommentOneDimensional(commentModel.ID)
		list = append(list, _list...)
	}
	return
}
