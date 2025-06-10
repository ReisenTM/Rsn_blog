package model

import (
	"blogX_server/global"
	"blogX_server/model/ctype"
	"blogX_server/model/enum"
	_ "embed"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ArticleModel 文章表
type ArticleModel struct {
	Model
	Title         string             `gorm:"size:32" json:"title" binding:"required"`
	Content       string             `json:"content" binding:"required"`
	Preview       string             `gorm:"size:256" json:"preview"`
	CategoryID    *uint              `json:"category_id"` //分类ID
	CategoryModel *CategoryModel     `gorm:"foreignKey:CategoryID" json:"-"`
	Tags          ctype.List         `gorm:"type:text" json:"tags"` //文章标签
	Cover         string             `gorm:"size:256" json:"cover"`
	UserID        uint               `json:"user_id"`
	UserModel     UserModel          `gorm:"foreignKey:UserID" json:"-"`
	ViewsCount    int                `json:"views_count"`
	FavorCount    int                `json:"favor_count"`
	CommentCount  int                `json:"comment_count"`
	CollectCount  int                `json:"collect_count"`
	OpenComment   bool               `json:"open_comment"` //开放评论
	Status        enum.ArticleStatus `json:"status"`       //状态:草稿，审核中，已发布
}

//go:embed mappings/article_mapping.json
var articleMapping string

func (ArticleModel) Mapping() string {
	return articleMapping
}

func (ArticleModel) Index() string {
	return "article_index"
}

// BeforeDelete 删除关联表记录
func (a *ArticleModel) BeforeDelete(tx *gorm.DB) (err error) {
	// 评论
	var commentList []CommentModel
	global.DB.Find(&commentList, "article_id = ?", a.ID).Delete(&commentList)
	// 点赞
	var favorList []UserArticleFavorModel
	global.DB.Find(&favorList, "article_id = ?", a.ID).Delete(&favorList)
	// 收藏
	var collectList []UserArticleCollectModel
	global.DB.Find(&collectList, "article_id = ?", a.ID).Delete(&collectList)
	// 置顶
	var topList []UserTopArticleModel
	global.DB.Find(&topList, "article_id = ?", a.ID).Delete(&topList)
	// 浏览
	var lookList []UserArticleHistoryModel
	global.DB.Find(&lookList, "article_id = ?", a.ID).Delete(&lookList)

	logrus.Infof("删除关联评论 %d 条", len(commentList))
	logrus.Infof("删除关联点赞 %d 条", len(favorList))
	logrus.Infof("删除关联收藏 %d 条", len(collectList))
	logrus.Infof("删除关联置顶 %d 条", len(topList))
	logrus.Infof("删除关联浏览 %d 条", len(lookList))
	return
}
