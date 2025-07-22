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

type CollectCreateRequest struct {
	Title string `json:"title" binding:"required,max=32"`
	ID    uint   `json:"id"`
	Info  string `json:"info"`
	Cover string `json:"cover"`
}

// CollectCreateView 创建或修改收藏夹
// 修改需要指定ID
func (ArticleApi) CollectCreateView(c *gin.Context) {
	cr := middleware.GetBind[CollectCreateRequest](c)
	//获取用户信息
	claim := jwts.GetClaims(c)

	var collect model.CollectModel
	if cr.ID == 0 {
		//创建
		err := global.DB.Take(&collect, "title = ? and user_id = ?", cr.Title, claim.UserID).Error
		if err == nil {
			resp.FailWithMsg("收藏夹名称重复", c)
			return
		}
		collect = model.CollectModel{
			Title:  cr.Title,
			UserID: claim.UserID,
			Info:   cr.Info,
			Cover:  cr.Cover,
		}
		err = global.DB.Create(&collect).Error
		if err != nil {
			resp.OKWithMsg("创建失败", c)
			logrus.Errorf("收藏夹 Create err:%v", err)
			return
		}
		resp.OKWithMsg("创建收藏夹成功", c)
		return
	}
	//修改
	err := global.DB.Take(&collect, "id = ? and user_id = ?", cr.ID, claim.UserID).Error
	if err != nil {
		resp.FailWithMsg("收藏夹不存在", c)
		return
	}
	err = global.DB.Updates(map[string]any{
		"title": cr.Title,
		"info":  cr.Info,
		"cover": cr.Cover,
	}).Error
	if err != nil {
		resp.FailWithMsg("修改失败", c)
		logrus.Errorf("收藏夹 Update err:%v", err)
		return
	}
	resp.OKWithMsg("修改成功", c)
}

type CollectListRequest struct {
	common.PageInfo
	UserID    uint `form:"user_id"`
	Type      int8 `form:"type" binding:"required,oneof=1 2 3"` // 1 查自己 2 查别人  3 后台
	ArticleID uint `form:"article_id"`
}
type CollectListResponse struct {
	model.CollectModel
	ArticleCount int    `json:"article_count"`
	Nickname     string `json:"nickname,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	ArticleUse   bool   `json:"articleUse,omitempty"` //是否已被收藏
}

func (ArticleApi) CollectListView(c *gin.Context) {
	cr := middleware.GetBind[CollectListRequest](c)
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

	_list, count, _ := common.ListQuery(model.CollectModel{
		UserID: cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title"},
		Preloads: preload,
	})
	var list = make([]CollectListResponse, 0)
	for _, i2 := range _list {
		item := CollectListResponse{
			CollectModel: i2,
			ArticleCount: len(i2.ArticleList),
			Nickname:     i2.UserModel.Nickname,
			Avatar:       i2.UserModel.Avatar,
		}
		for _, model := range i2.ArticleList {
			if model.ArticleID == cr.ArticleID {
				item.ArticleUse = true
				break
			}
		}
		list = append(list, item)
	}
	resp.OkWithList(list, count, c)
}

func (ArticleApi) CollectRemoveView(c *gin.Context) {
	cr := middleware.GetBind[model.RemoveRequest](c)
	var list []model.CollectModel
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
			resp.FailWithMsg("删除收藏夹失败", c)
			return
		}
	}

	msg := fmt.Sprintf("删除收藏夹成功 共删除%d条", len(list))

	resp.OKWithMsg(msg, c)
}
