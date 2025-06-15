package article_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ArticleCategoryRequest struct {
	Title string `json:"title"`
	ID    uint   `json:"id"`
}

// CategoryCreateView 创建或修改分类
func (ArticleApi) CategoryCreateView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCategoryRequest](c)
	//获取用户信息
	claim := jwts.GetClaims(c)

	var category model.CategoryModel
	if cr.ID == 0 {
		//创建
		err := global.DB.Take(&category, "title = ? and user_id = ?", cr.Title, claim.UserID).Error
		if err == nil {
			resp.FailWithMsg("分类名称重复", c)
			return
		}
		category = model.CategoryModel{
			Title:  cr.Title,
			UserID: claim.UserID,
		}
		err = global.DB.Create(&category).Error
		if err != nil {
			resp.OKWithMsg("创建失败", c)
			logrus.Errorf("Category Create err:%v", err)
			return
		}
		resp.OKWithMsg("创建分类成功", c)
		return
	}
	//修改
	err := global.DB.Take(&category, "id = ? and user_id = ?", cr.ID, claim.UserID).Error
	if err != nil {
		resp.FailWithMsg("分类不存在", c)
		return
	}
	err = global.DB.Updates(&category).Error
	if err != nil {
		resp.FailWithMsg("修改失败", c)
		logrus.Errorf("Category Update err:%v", err)
		return
	}
	resp.OKWithMsg("修改成功", c)
}

type ArticleCategoryListRequest struct {
	common.PageInfo
	UserID uint `form:"user_id"`
	Type   int8 `form:"type" binding:"required,oneof=1 2 3"` // 1 查自己 2 查别人  3 后台
}
type CategoryListResponse struct {
	model.CategoryModel
	ArticleCount int    `json:"articleCount"`
	Nickname     string `json:"nickname,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
}

func (ArticleApi) CategoryListView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCategoryListRequest](c)
	var preload = []string{"ArticleList"}

	switch cr.Type {
	case 1:
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			resp.FailWithError(err, c)
			return
		}
		cr.UserID = claims.UserID
	case 2:
		if cr.UserID == 0 {
			resp.FailWithMsg("必须指定目标用户", c)
			return
		}
	case 3:
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			resp.FailWithError(err, c)
			return
		}
		if claims.Role != enum.RoleAdminType {
			resp.FailWithMsg("权限错误", c)
			return
		}
		preload = append(preload, "UserModel")
	}

	_list, count, _ := common.ListQuery(model.CategoryModel{
		UserID: cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title"},
		Preloads: preload,
	})
	var list = make([]CategoryListResponse, 0)
	for _, v := range _list {
		list = append(list, CategoryListResponse{
			CategoryModel: v,
			ArticleCount:  len(v.ArticleList), //每个分类对应的文章数量，不是分类数
			Nickname:      v.UserModel.Nickname,
			Avatar:        v.UserModel.Avatar,
		})
	}
	resp.OkWithList(list, count, c)
}

func (ArticleApi) CategoryRemoveView(c *gin.Context) {
	cr := middleware.GetBind[model.RemoveRequest](c)
	var list []model.CategoryModel
	query := global.DB.Where("id in ?", cr.IDList)
	claims := jwts.GetClaims(c)
	if claims.Role != enum.RoleAdminType {
		//不是管理员只能删自己的
		query.Where("user_id = ?", claims.UserID)
	}

	global.DB.Where(query).Find(&list)

	if len(list) > 0 {
		err := global.DB.Delete(&list).Error
		if err != nil {
			resp.FailWithMsg("删除分类失败", c)
			return
		}
	}

	msg := fmt.Sprintf("删除分类成功 共删除%d条", len(list))

	resp.OKWithMsg(msg, c)
}

// CategoryOptionsView 分类列表，显示用户的所有分类，方便前端调用
func (ArticleApi) CategoryOptionsView(c *gin.Context) {
	claims := jwts.GetClaims(c)

	var list []model.OptionsResponse[uint]
	global.DB.Model(model.CategoryModel{}).Where("user_id = ?", claims.UserID).
		Select("id as value", "title as label").Scan(&list)

	resp.OkWithData(list, c)

}
