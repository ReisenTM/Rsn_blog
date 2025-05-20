package banner_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

type BannerApi struct{}
type BannerCreateRequest struct {
	Cover string `json:"cover" binding:"required"`
	Href  string `json:"href"` //可能不需要点击跳转
	Show  bool   `json:"show"`
}

// BannerCreateView 封面上传
func (BannerApi) BannerCreateView(c *gin.Context) {
	var req BannerCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	err = global.DB.Create(&model.BannerModel{
		Cover: req.Cover,
		Href:  req.Href,
		Show:  req.Show,
	}).Error
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	resp.OKWithMsg("封面上传成功", c)
}

type BannerListRequest struct {
	common.PageInfo
	Show bool `form:"show"`
}

// BannerListCreateView 列表请求
func (BannerApi) BannerListCreateView(c *gin.Context) {
	var st BannerListRequest
	err := c.ShouldBindQuery(&st)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	var list []model.BannerModel
	list, count, err := common.ListQuery(model.BannerModel{
		Show: st.Show,
	}, common.Options{
		PageInfo: st.PageInfo,
	})
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	resp.OkWithList(list, count, c)
}

// BannerRemoveView 删除封面
func (BannerApi) BannerRemoveView(c *gin.Context) {
	var rm model.DeleteRequest
	err := c.ShouldBindJSON(&rm)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	var list []model.BannerModel
	err = global.DB.Find(&list, "id in ?", rm.IDList).Error
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	global.DB.Delete(&list)
	resp.OKWithMsg(fmt.Sprintf("欲删除%d条数据,成功删除%d条", len(rm.IDList), len(list)), c)
}

// BannerUpdateView 修改封面
func (BannerApi) BannerUpdateView(c *gin.Context) {
	var id model.IDRequest
	err := c.ShouldBindUri(&id)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	var cr BannerCreateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	var m model.BannerModel
	err = global.DB.Take(&m, id).Error
	if err != nil {
		resp.FailWithMsg("不存在的banner", c)
		return
	}

	err = global.DB.Model(&m).Updates(map[string]any{
		"cover": cr.Cover,
		"href":  cr.Href,
		"show":  cr.Show,
	}).Error
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	resp.OKWithMsg("banner更新成功", c)
}
