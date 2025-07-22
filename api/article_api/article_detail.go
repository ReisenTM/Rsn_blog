package article_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/service/redis_service/redis_article"
	"blogX_server/utils/jwts"
	"github.com/gin-gonic/gin"
)

type ArticleDetailResponse struct {
	model.ArticleModel
	Username      string  `json:"username"`   //用户名
	Nickname      string  `json:"nickname"`   //用户昵称
	UserAvatar    string  `json:"avatar"`     //用户头像
	IsFavor       bool    `json:"is_favor"`   //是否点赞
	IsCollcet     bool    `json:"is_collcet"` //是否收藏
	CategoryTitle *string `json:"category_title"`
}

func (ArticleApi) ArticleDetailView(c *gin.Context) {
	cr := middleware.GetBind[model.IDRequest](c)
	// 未登录的用户，只能看到发布成功的文章

	// 登录用户，能看到自己的所有文章

	// 管理员，能看到全部的文章
	var article model.ArticleModel
	err := global.DB.Preload("UserModel").Preload("CategoryModel").Take(&article, cr.ID).Error
	if err != nil {
		resp.FailWithMsg("不存在的文章", c)
		return
	}

	claim, err := jwts.ParseTokenByGin(c)
	if err != nil {
		// 没登录的
		if article.Status != enum.ArticleStatusPublished {
			resp.FailWithMsg("文章不存在", c)
			return
		}
	}
	var res = ArticleDetailResponse{
		ArticleModel: article,
		Username:     article.UserModel.Username,
		Nickname:     article.UserModel.Nickname,
		UserAvatar:   article.UserModel.Avatar,
	}
	//如果用户登录了
	if claim != nil && err == nil {
		if claim.Role == enum.RoleUserType {
			if claim.UserID != article.UserID {
				//普通用户查看不是自己的文章
				if article.Status != enum.ArticleStatusPublished {
					//只能看已经发布的文章
					resp.FailWithMsg("文章不存在", c)
					return
				}
			}
		}
		//看是否点赞
		var fm model.UserArticleFavorModel
		err = global.DB.Take(&fm, "user_id = ? AND article_id = ?", article.UserID, article.ID).Error
		if err == nil {
			res.IsFavor = true
		}
		//看是否加入了收藏夹
		var cm model.UserArticleCollectModel
		err = global.DB.Take(&cm, "user_id = ? AND article_id = ?", article.UserID, article.ID).Error
		if err == nil {
			res.IsCollcet = true
		}
	}
	favorCount := redis_article.GetCacheFavor(cr.ID)
	viewCount := redis_article.GetCacheLook(cr.ID)
	commentCount := redis_article.GetCacheComment(cr.ID)
	collectCount := redis_article.GetCacheCollect(cr.ID)
	res.FavorCount = favorCount + article.FavorCount
	res.ViewsCount = viewCount + article.ViewsCount
	res.CommentCount = commentCount + article.CommentCount
	res.CollectCount = collectCount + article.CollectCount
	if article.CategoryModel != nil {
		res.CategoryTitle = &article.CategoryModel.Title
	}
	resp.OkWithData(res, c)
}
