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
	"blogX_server/utils/sql"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ArticleListRequest struct {
	common.PageInfo
	UserID     uint               `form:"user_id"`
	Type       int8               `form:"type" binding:"required,oneof=1 2 3"` //用户查别人，查自己，管理员查别人
	CategoryID *uint              `form:"category_id"`
	Status     enum.ArticleStatus `form:"status"`
	CollectID  uint               `form:"collect"`
}
type ArticleListResponse struct {
	model.ArticleModel
	UserTop       bool    `json:"user_top"`       // 是否是用户置顶
	AdminTop      bool    `json:"admin_top"`      // 是否是管理员置顶
	CategoryTitle *string `json:"category_title"` //分类
	Avatar        string  `json:"avatar"`         //用户头像
	Nickname      string  `json:"nickname"`       //用户昵称
}

// ArticleListView 文章列表
func (ArticleApi) ArticleListView(c *gin.Context) {
	cr := middleware.GetBind[ArticleListRequest](c)
	switch cr.Type {
	case 1:
		//如果是用户查别人
		// 查别人。用户id就是必填的
		if cr.UserID == 0 {
			resp.FailWithMsg("用户id必填", c)
			return
		}
		if cr.Page > 2 || cr.Limit > 10 {
			resp.FailWithMsg("查询更多，请登录", c)
			return
		}
		//只能查已发布的
		cr.Status = enum.ArticleStatusPublished
		if cr.CollectID != 0 {
			//如果根据查收藏夹内文章
			var userConf model.UserConfModel
			err := global.DB.Take(&userConf, "user_id = ?", cr.UserID).Error
			if err != nil {
				resp.FailWithMsg("用户不存在", c)
				return
			}
			//如果用户设置了不公开
			if !userConf.OpenCollection {
				resp.FailWithMsg("用户未开启我的收藏", c)
				return
			}
		}
	case 2:
		// 查自己的
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			resp.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.UserID
	case 3:
		// 管理员
		claims, err := jwts.ParseTokenByGin(c)
		if !(err == nil && claims.Role == enum.RoleAdminType) {
			resp.FailWithMsg("角色错误", c)
			return
		}
	}
	//预设排序映射
	var OrderMap = map[string]bool{
		"views_count desc":   true, //降序
		"favor_count desc":   true,
		"comment_count desc": true,
		"collect_count desc": true,
		"views_count asc":    true, //升序
		"favor_count asc":    true,
		"comment_count asc":  true,
		"collect_count asc":  true,
	}
	if cr.Order != "" {
		_, ok := OrderMap[cr.Order]
		if !ok {
			resp.FailWithMsg("不支持的排序方式", c)
			return
		}
	}

	//找是否有置顶文章
	var adminTopMap = make(map[uint]bool) //管理员置顶映射 id->y/n
	var userTopMap = make(map[uint]bool)  //用户置顶映射
	var topArticleIDList []uint           //置顶文章id list
	var userTopQuery = global.DB.Where("")
	if cr.UserID != 0 {
		userTopQuery.Where("user_id = ?", cr.UserID)
	}
	var userTopArticleList []model.UserTopArticleModel
	//按照创建时间倒序寻找
	global.DB.Preload("UserModel").Order("created_at desc").Where(userTopQuery).Find(&userTopArticleList)
	for _, tm := range userTopArticleList {
		topArticleIDList = append(topArticleIDList, tm.ArticleID)
		if tm.UserModel.Role == enum.RoleAdminType {
			//如果是管理员置顶
			adminTopMap[tm.ArticleID] = true
		}
		userTopMap[tm.ArticleID] = true
	}
	//列表查询配置
	var options = common.Options{
		PageInfo:     cr.PageInfo,
		Likes:        []string{"title"},
		DefaultOrder: "created_at desc",
		Preloads:     []string{"UserModel", "CategoryModel"},
	}
	//如果存在置顶文章
	if len(topArticleIDList) > 0 {
		//优先查看置顶文章，其余文章倒序
		options.DefaultOrder = fmt.Sprintf("%s, created_at desc", sql.ConvertSliceOrderSql(topArticleIDList))
	}
	//执行列表查询
	_list, count, _ := common.ListQuery(model.ArticleModel{
		UserID:     cr.UserID,
		CategoryID: cr.CategoryID,
		Status:     cr.Status,
	}, options)
	var list = make([]ArticleListResponse, 0)
	collectMap := redis_article.GetAllCacheCollect()
	viewMap := redis_article.GetAllCacheLook()
	favorMap := redis_article.GetAllCacheFavor()
	commentMap := redis_article.GetAllCacheComment()

	for _, mod := range _list {
		//不需要显示正文
		mod.Content = ""
		mod.FavorCount = mod.FavorCount + favorMap[mod.ID]
		mod.CollectCount = mod.CollectCount + collectMap[mod.ID]
		mod.ViewsCount = mod.ViewsCount + viewMap[mod.ID]
		mod.CommentCount = mod.CommentCount + commentMap[mod.ID]
		data := ArticleListResponse{
			ArticleModel: mod,
			AdminTop:     adminTopMap[mod.ID], //看文章是否在置顶里
			UserTop:      userTopMap[mod.ID],
			Avatar:       mod.UserModel.Avatar,
			Nickname:     mod.UserModel.Nickname,
		}

		if mod.CategoryModel != nil {
			//因为是指针
			data.CategoryTitle = &mod.CategoryModel.Title
		}
		list = append(list, data)
	}

	resp.OkWithList(list, count, c)
}
