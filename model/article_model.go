package model

import (
	"blogX_server/global"
	"blogX_server/model/ctype"
	"blogX_server/model/enum"
	"blogX_server/service/text_service"
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
func (a *ArticleModel) AfterCreate(tx *gorm.DB) (err error) {
	// 创建文章之后的钩子函数
	// 只有发布中的文章会放到全文搜索里面去
	if a.Status != enum.ArticleStatusPublished {
		return nil
	}
	textList := text_service.MdContentTransformation(a.ID, a.Title, a.Content)
	var list []TextModel
	if len(textList) == 0 {
		return nil
	}
	//手动控制只设置关键字段，避免某些默认字段（如 ID、CreatedAt）被错误赋值
	for _, model := range textList {
		list = append(list, TextModel{
			ArticleID: model.ArticleID,
			Head:      model.Head,
			Body:      model.Body,
		})
	}
	err = tx.Create(&list).Error
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return nil
}

func (a *ArticleModel) AfterDelete(tx *gorm.DB) (err error) {
	// 删除之后
	var textList []TextModel
	tx.Find(&textList, "article_id = ?", a.ID)
	if len(textList) > 0 {
		logrus.Infof("删除全文记录 %d", len(textList))
		tx.Delete(&textList)
	}
	return nil
}

func (a *ArticleModel) AfterUpdate(tx *gorm.DB) (err error) {
	// 正文发生了变化，才做转换
	a.AfterDelete(tx)
	a.AfterCreate(tx)
	return nil
}
