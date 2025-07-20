package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserArticleTopRequest struct {
	ArticleID uint `json:"article_id" binding:"required"`
	Type      uint `json:"type" binding:"required,oneof=1 2"` //1 用户置顶 2 管理员置顶
}

// UserArticleTopView 用户文章置顶
func (UserApi) UserArticleTopView(c *gin.Context) {
	cr := middleware.GetBind[UserArticleTopRequest](c)
	claims := jwts.GetClaims(c)
	//先看文章存在不存在
	var article model.ArticleModel
	err := global.DB.Take(&article, cr.ArticleID).Error
	if err != nil {
		resp.FailWithMsg("文章不存在", c)
		return
	}
	//再看置顶类型
	switch cr.Type {
	case 1:
		// 用户置顶文章
		// 验证文章是不是自己的，并且是已发布的
		if article.UserID != claims.UserID {
			resp.FailWithMsg("用户只能置顶自己的文章", c)
			return
		}

		if article.Status != enum.ArticleStatusPublished {
			resp.FailWithMsg("用户只能置顶已发布的文章", c)
			return
		}

		// 判断之前自己有没有置顶过
		var userTopArticleList []model.UserTopArticleModel
		global.DB.Find(&userTopArticleList, "user_id = ?",
			claims.UserID)
		// 查不到  自己从来没有置顶过文章
		if len(userTopArticleList) == 0 {
			//置顶
			err = global.DB.Create(&model.UserTopArticleModel{
				UserID:    claims.UserID,
				ArticleID: article.ID,
			}).Error
			if err != nil {
				logrus.Errorf("create user top article err: %v", err)
				resp.FailWithError(err, c)
			}
			resp.OKWithMsg("置顶成功", c)
			return
		}
		if len(userTopArticleList) == 1 {
			uta := userTopArticleList[0]
			if uta.ArticleID != cr.ArticleID {
				resp.FailWithMsg("普通用户只能置顶一篇文章", c)
				return
			}
		}
		//如果对已经置顶的文章操作，就是取消置顶
		uta := userTopArticleList[0]
		global.DB.Delete(&uta, "article_id = ?", uta.ArticleID)
		resp.OKWithMsg("取消置顶成功", c)
	case 2:
		//如果是管理员置顶
		//置顶
		err = global.DB.Create(&model.UserTopArticleModel{
			UserID:    claims.UserID,
			ArticleID: article.ID,
		}).Error
		if err != nil {
			logrus.Errorf("create user top article err: %v", err)
			resp.FailWithError(err, c)
		}
		resp.OKWithMsg("置顶成功", c)

	}
}
