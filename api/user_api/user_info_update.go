package user_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum"
	"blogX_server/utils/jwts"
	"blogX_server/utils/mps"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// 指针是为了实现原位修改，可以实现只对被修改的数据进行更新
type UserInfoUpdateRequest struct {
	Username    *string   `json:"username" s-u:"username"`
	Nickname    *string   `json:"nickname" s-u:"nickname"`
	Avatar      *string   `json:"avatar" s-u:"avatar"`
	Profile     *string   `json:"profile" s-u:"profile"`
	LikeTags    *[]string `json:"like_tags" s-u-c:"like_tags"`
	OpenCollect *bool     `json:"open_collect" s-u-c:"open_collect"`   // 公开我的收藏
	OpenFollow  *bool     `json:"open_follow" s-u-c:"open_follow"`     // 公开我的关注
	OpenFans    *bool     `json:"open_fans" s-u-c:"open_fans"`         // 公开我的粉丝
	HomeStyleID *uint     `json:"home_style_id" s-u-c:"home_style_id"` // 主页样式的id
}

// UserInfoUpdateView 个人信息修改
func (UserApi) UserInfoUpdateView(c *gin.Context) {
	var cr UserInfoUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		resp.FailWithError(err, c)
		return
	}
	userMap := mps.StructToMap(cr, "s-u")
	userConfMap := mps.StructToMap(cr, "s-u-c")
	fmt.Println("userMap", userMap)
	fmt.Println("userConfMap", userConfMap)

	claims := jwts.GetClaims(c)

	if len(userMap) > 0 {
		//更新用户信息
		var userModel model.UserModel
		err = global.DB.Preload("UserConfModel").Take(&userModel, claims.UserID).Error
		if err != nil {
			resp.FailWithMsg("用户不存在", c)
			return
		}
		// 判断
		if cr.Username != nil {
			var userCount int64
			global.DB.Debug().Model(model.UserModel{}).
				Where("username = ? and id <> ?", *cr.Username, claims.UserID).
				Count(&userCount)
			fmt.Println(*cr.Username, userCount)
			if userCount > 0 {
				resp.FailWithMsg("该用户名被使用", c)
				return
			}
			if *cr.Username != userModel.Username {
				// 如果和我的用户名是不一样的
				var uud = userModel.UserConfModel.UpdateMark
				if uud != nil {
					if time.Now().Sub(*uud).Hours() < 720 {
						resp.FailWithMsg("用户名30天内只能修改一次", c)
						return
					}
				}
				userConfMap["update_mark"] = time.Now()
			}
		}

		if cr.Nickname != nil || cr.Avatar != nil {
			if userModel.RegSource == enum.RegisterQQSourceType {
				resp.FailWithMsg("QQ注册的用户不能修改昵称和头像", c)
				return
			}
		}

		err = global.DB.Model(&userModel).Updates(userMap).Error
		if err != nil {
			resp.FailWithMsg("用户信息修改失败", c)
			return
		}
	}
	if len(userConfMap) > 0 {
		//更新用户配置信息
		var userConfModel model.UserConfModel
		err = global.DB.Take(&userConfModel, "user_id = ?", claims.UserID).Error
		if err != nil {
			resp.FailWithMsg("用户配置信息不存在", c)
			return
		}
		err = global.DB.Model(&userConfModel).Updates(userConfMap).Error
		if err != nil {
			resp.FailWithMsg("用户信息修改失败", c)
			return
		}
	}

	resp.OKWithMsg("用户信息修改成功", c)

}
