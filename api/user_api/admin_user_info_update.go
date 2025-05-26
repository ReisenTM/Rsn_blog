package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/utils/mps"
	"github.com/gin-gonic/gin"
)

type AdminUserInfoUpdateRequest struct {
	UserID   uint           `json:"user_id" binding:"required"`
	Username *string        `json:"username" s-u:"username"`
	Nickname *string        `json:"nickname" s-u:"nickname"`
	Avatar   *string        `json:"avatar" s-u:"avatar"`
	Profile  *string        `json:"profile" s-u:"profile"`
	Role     *enum.RoleType `json:"role" s-u:"role"`
}

// AdminUserInfoUpdateView 个人信息修改
func (UserApi) AdminUserInfoUpdateView(c *gin.Context) {
	var cr AdminUserInfoUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	userMap := mps.StructToMap(cr, "s-u")
	var user model.UserModel
	err = global.DB.Take(&user, cr.UserID).Error
	if err != nil {
		resp.FailWithMsg("用户不存在", c)
		return
	}

	err = global.DB.Model(&user).Updates(userMap).Error
	if err != nil {
		resp.FailWithMsg("用户信息修改失败", c)
		return
	}

	resp.OKWithMsg("用户信息修改成功", c)
}
