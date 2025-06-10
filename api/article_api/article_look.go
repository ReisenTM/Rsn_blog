package article_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/service/redis_service/redis_article"
	"blogX_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type ArticleLookRequest struct {
	ArticleID uint `json:"article_id" binding:"required"`
	TimeTake  int  `json:"time_take"` // 读文章一共用了多久
}

// ArticleLookView 浏览
func (ArticleApi) ArticleLookView(c *gin.Context) {
	cr := middleware.GetBind[ArticleLookRequest](c)
	// TODO: 未登录用户，浏览量如何加

	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		resp.OKWithMsg("未登录", c)
		return
	}
	// 引入缓存
	// 当天这个用户请求这个文章之后，将用户id和文章id作为key存入缓存，在这里进行判断，如果存在就直接返回
	if redis_article.GetUserArticleHistoryCache(cr.ArticleID, claims.UserID) {
		resp.OKWithMsg("成功", c)
		return
	}
	//文章存在吗
	var article model.ArticleModel
	err = global.DB.Take(&article, "status = ? and id = ?", enum.ArticleStatusPublished, cr.ArticleID).Error
	if err != nil {
		resp.FailWithMsg("文章不存在", c)
		return
	}
	// 查这个文章今天有没有在足迹里面
	var history model.UserArticleHistoryModel
	err = global.DB.Take(&history,
		"user_id = ? and article_id = ? and created_at < ? and created_at > ?",
		claims.UserID, cr.ArticleID,
		time.Now().Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02")+" 00:00:00",
	).Error
	if err == nil {
		resp.OKWithMsg("成功", c)
		return
	}
	//不在则创建足迹
	err = global.DB.Create(&model.UserArticleHistoryModel{
		UserID:    claims.UserID,
		ArticleID: cr.ArticleID,
	}).Error
	if err != nil {
		resp.FailWithMsg("失败", c)
		return
	}
	//浏览量+1
	redis_article.SetCacheLook(cr.ArticleID, true)
	//足迹添加
	redis_article.SetUserArticleHistoryCache(cr.ArticleID, claims.UserID)

	resp.OKWithMsg("成功", c)
}

type ArticleLookListRequest struct {
	common.PageInfo
	UserID uint `form:"user_id"`
	Type   int8 `form:"type" binding:"required,oneof=1 2"` // 1 查自己 2 查别人
}
type ArticleLookListResponse struct {
	ID        uint      `json:"id"`       // 浏览记录的id
	LookDate  time.Time `json:"lookDate"` // 浏览的时间
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	UserID    uint      `json:"userID"`
	ArticleID uint      `json:"articleID"`
}

// ArticleLookListView 文章浏览列表
func (ArticleApi) ArticleLookListView(c *gin.Context) {
	cr := middleware.GetBind[ArticleLookListRequest](c)
	claims := jwts.GetClaims(c)

	switch cr.Type {
	case 1:
		cr.UserID = claims.UserID
	}

	_list, count, _ := common.ListQuery(model.UserArticleHistoryModel{
		UserID: cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"UserModel", "ArticleModel"},
	})

	var list = make([]ArticleLookListResponse, 0)
	for _, historyModel := range _list {
		list = append(list, ArticleLookListResponse{
			ID:        historyModel.ID,
			LookDate:  historyModel.CreatedAt,
			Title:     historyModel.ArticleModel.Title,
			Cover:     historyModel.ArticleModel.Cover,
			Nickname:  historyModel.UserModel.Nickname,
			Avatar:    historyModel.UserModel.Avatar,
			UserID:    historyModel.UserID,
			ArticleID: historyModel.ArticleID,
		})
	}

	resp.OkWithList(list, count, c)

}

// ArticleLookRemoveView 删除历史记录
func (ArticleApi) ArticleLookRemoveView(c *gin.Context) {
	cr := middleware.GetBind[model.RemoveRequest](c)

	claims := jwts.GetClaims(c)
	var list []model.UserArticleHistoryModel
	global.DB.Find(&list, "user_id = ? and id in ?", claims.UserID, cr.IDList)

	if len(list) > 0 {
		err := global.DB.Delete(&list).Error
		if err != nil {
			resp.FailWithMsg("足迹删除失败", c)
			return
		}
	}

	resp.OKWithMsg(fmt.Sprintf("删除足迹成功 共删除%d条", len(list)), c)
}
