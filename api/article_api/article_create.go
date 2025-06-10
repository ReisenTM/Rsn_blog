package article_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/ctype"
	"blogX_server/model/enum"
	"blogX_server/utils/jwts"
	"blogX_server/utils/markdown"
	"blogX_server/utils/xss"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ArticleCreateRequest struct {
	Title       string             `json:"title" binding:"required"`
	Content     string             `json:"content" binding:"required"`
	Preview     string             `json:"preview"`
	CategoryID  *uint              `json:"category_id"`
	Tags        ctype.List         `json:"tags" binding:"required"`
	OpenComment bool               `json:"open_comment"`
	Status      enum.ArticleStatus `json:"status"`
	Cover       string             `json:"cover"`
}

// CreateArticleView 创建文章
func (ArticleApi) CreateArticleView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCreateRequest](c)
	//获取用户id
	user, err := jwts.GetClaims(c).GetUser()
	if err != nil {
		resp.FailWithMsg("获取用户信息失败", c)
		return
	}
	if global.Config.Site.SiteInfo.Mode == 2 {
		if user.Role != enum.RoleAdminType {
			resp.FailWithMsg("博客模式下，普通用户不能发文章", c)
			return
		}
	}
	//判断分类id是否是自己创建
	var categoryModel model.CategoryModel
	if cr.CategoryID != nil {
		err = global.DB.First(&categoryModel, "id = ? and user_id = ?",
			*cr.CategoryID, user.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.FailWithMsg("文章分类不存在", c)
			return
		}
	}
	//正文防止xss注入
	cr.Content = xss.XssFilter(cr.Content)
	//如果用户没主动写简介，取正文前100字符生成简介
	if cr.Preview == "" {
		preview, err := markdown.GetPreviewContent(cr.Content, 100)
		if err != nil {
			logrus.Errorf("GetPreviewContent err: %v", err)
			resp.FailWithMsg("正文解析错误", c)
			return
		}
		cr.Preview = preview
	}
	// 正文内容图片转存
	// 1.图片过多，同步做，接口耗时高  异步做，

	var m = model.ArticleModel{
		Title:       cr.Title,
		Content:     cr.Content,
		Preview:     cr.Preview,
		CategoryID:  cr.CategoryID,
		Tags:        cr.Tags,
		Cover:       cr.Cover,
		UserID:      user.ID,
		OpenComment: cr.OpenComment,
		Status:      cr.Status,
	}
	err = global.DB.Create(&m).Error
	if err != nil {
		resp.FailWithError(err, c)
		logrus.Errorf("创建文章-保存到数据库失败")
		return
	}
	resp.OKWithMsg("创建文章成功", c)
}
