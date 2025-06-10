package image_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/service/log_service"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ImageApi struct {
}
type ImageListResponse struct {
	model.ImageModel
	WebPath string `json:"webPath"`
}

// ImageListView 列表请求图片
func (ImageApi) ImageListView(c *gin.Context) {
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}

	_list, count, _ := common.ListQuery(model.ImageModel{}, common.Options{
		PageInfo: cr,
		Likes:    []string{"filename"},
	})
	var list = make([]ImageListResponse, 0)
	for _, m := range _list {
		list = append(list, ImageListResponse{
			ImageModel: m,
			WebPath:    m.NetPath(),
		})
	}
	resp.OkWithList(list, count, c)
}

// ImageRemoveView 通过id列表删除对应记录
func (ImageApi) ImageRemoveView(c *gin.Context) {
	var cr model.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	log := log_service.GetLog(c)
	log.SetShowReq()
	log.SetShowRes()

	var list []model.ImageModel
	global.DB.Find(&list, "id in ?", cr.IDList)

	var successCount, errCount int64
	if len(list) > 0 {
		successCount = global.DB.Delete(&list).RowsAffected
	}
	errCount = int64(len(list)) - successCount

	msg := fmt.Sprintf("操作成功，成功%d 失败%d", successCount, errCount)

	resp.OKWithMsg(msg, c)
}
