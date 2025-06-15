package model

import "blogX_server/model/enum/message_type_enum"

type MessageModel struct {
	Model
	Type               message_type_enum.Type `json:"type"`
	RevUserID          uint                   `json:"rev_user_id"` // 接收人的id
	ActionUserID       uint                   `json:"action_user_id"`
	ActionUserNickname string                 `json:"action_user_nickname"`
	ActionUserAvatar   string                 `json:"action_user_avatar"`
	Title              string                 `json:"title"`
	Content            string                 `json:"content"`
	ArticleID          uint                   `json:"article_id"`
	ArticleTitle       string                 `json:"article_title"`
	CommentID          uint                   `json:"comment_id"`
	LinkTitle          string                 `json:"link_title"`
	LinkHref           string                 `json:"link_href"`
	IsRead             bool                   `json:"is_read"`
}
