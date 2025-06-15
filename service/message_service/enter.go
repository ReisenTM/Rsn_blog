package message_service

import (
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum/message_type_enum"
	"github.com/sirupsen/logrus"
)

// InsertCommentMessage 插入一条评论消息
func InsertCommentMessage(commentModel model.CommentModel) {
	global.DB.Preload("UserModel").Preload("ArticleModel").Take(&commentModel)
	err := global.DB.Create(&model.MessageModel{
		Type:               message_type_enum.CommentType,
		RevUserID:          commentModel.ArticleModel.UserID,
		ActionUserID:       commentModel.UserID,
		ActionUserNickname: commentModel.UserModel.Nickname,
		ActionUserAvatar:   commentModel.UserModel.Avatar,
		Content:            commentModel.Content,
		ArticleID:          commentModel.ArticleID,
		ArticleTitle:       commentModel.ArticleModel.Title,
		CommentID:          commentModel.ID,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertReplyMessage 插入一条回复消息
func InsertReplyMessage(commentModel model.CommentModel) {
	//TODO:自己回复自己怎么办?
	global.DB.Preload("ParentModel").Preload("UserModel").Preload("ArticleModel").Take(&commentModel)
	err := global.DB.Create(&model.MessageModel{
		Type:               message_type_enum.ReplyType,
		RevUserID:          commentModel.ParentModel.UserID, //注意需要的是父评论的UID而不是父评论本身的ID
		ActionUserID:       commentModel.UserID,
		ActionUserNickname: commentModel.UserModel.Nickname,
		ActionUserAvatar:   commentModel.UserModel.Avatar,
		Content:            commentModel.Content,
		ArticleID:          commentModel.ArticleID,
		ArticleTitle:       commentModel.ArticleModel.Title,
		CommentID:          commentModel.ID,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertFavorArticleMessage 点赞文章的消息
func InsertFavorArticleMessage(favorModel model.UserArticleFavorModel) {
	global.DB.Preload("UserModel").Preload("ArticleModel").Take(&favorModel)
	err := global.DB.Create(&model.MessageModel{
		Type:               message_type_enum.FavorArticleType,
		RevUserID:          favorModel.ArticleModel.UserID,
		ActionUserID:       favorModel.UserID,
		ActionUserNickname: favorModel.UserModel.Nickname,
		ActionUserAvatar:   favorModel.UserModel.Avatar,
		ArticleID:          favorModel.ArticleID,
		ArticleTitle:       favorModel.ArticleModel.Title,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertFavorCommentMessage 点赞评论的消息
func InsertFavorCommentMessage(m model.UserCommentFavorModel) {
	global.DB.Preload("CommentModel.ArticleModel").Preload("UserModel").Take(&m)
	err := global.DB.Create(&model.MessageModel{
		Type:               message_type_enum.FavorCommentType,
		RevUserID:          m.CommentModel.UserID,
		ActionUserID:       m.UserID,
		ActionUserNickname: m.UserModel.Nickname,
		ActionUserAvatar:   m.UserModel.Avatar,
		Content:            m.CommentModel.Content,
		ArticleID:          m.CommentModel.ArticleID,
		ArticleTitle:       m.CommentModel.ArticleModel.Title,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertCollectArticleMessage 收藏文章的消息
func InsertCollectArticleMessage(m model.UserArticleCollectModel) {
	global.DB.Preload("ArticleModel").Preload("UserModel").Take(&m)
	err := global.DB.Create(&model.MessageModel{
		Type:               message_type_enum.CollectArticleType,
		RevUserID:          m.ArticleModel.UserID,
		ActionUserID:       m.UserID,
		ActionUserNickname: m.UserModel.Nickname,
		ActionUserAvatar:   m.UserModel.Avatar,
		ArticleID:          m.ArticleID,
		ArticleTitle:       m.ArticleModel.Title,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}

// InsertSystemMessage 插入系统消息
func InsertSystemMessage(revUserID uint, title string, content string, linkTitle string, linkHref string) {
	err := global.DB.Create(&model.MessageModel{
		Type:      message_type_enum.SystemType,
		RevUserID: revUserID,
		Title:     title,
		Content:   content,
		LinkTitle: linkTitle,
		LinkHref:  linkHref,
	}).Error
	if err != nil {
		logrus.Error(err)
	}
}
