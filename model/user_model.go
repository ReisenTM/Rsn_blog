package model

import (
	"blogX_server/model/enum"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"reflect"
	"time"
)

// UserModel 用户表
type UserModel struct {
	Model
	Username      string                  `json:"user_name" gorm:"size=36"` //用户名
	Password      string                  `json:"password" gorm:"size=128"` //密码
	Nickname      string                  `json:"nick_name" gorm:"size=36"` //昵称
	Email         string                  `json:"email" gorm:"size=128"`    //邮箱
	Profile       string                  `json:"profile" gorm:"size:255"`  //简介
	RegSource     enum.RegisterSourceType `json:"reg_source"`               //注册来源
	Avatar        string                  `json:"avatar" gorm:"size:256"`   //头像(地址)
	OpenID        string                  `json:"open_id" gorm:"size:32"`   //第三方登录的唯一id
	Role          enum.RoleType           `json:"role" gorm:"size:4"`       //角色：1.用户2.管理员3.访客
	UserConfModel *UserConfModel          `gorm:"foreignKey:UserID"  json:"-"`
	IP            string                  `json:"ip"`
	Region        string                  `json:"region"` //ip归属地
	ArticleList   []ArticleModel          `gorm:"foreignKey:UserID" json:"-"`
	LoginList     []UserLoginModel        `gorm:"foreignKey:UserID" json:"-"`
}

// AfterCreate 随userModel一起创建
func (u *UserModel) AfterCreate(tx *gorm.DB) error {
	err := tx.Create(&UserConfModel{UserID: u.ID}).Error
	err = tx.Create(&UserMessageConfModel{UserID: u.ID, OpenCommentMessage: true, OpenFavorMessage: true, OpenPrivateChat: true}).Error
	return err
}
func (u *UserModel) BeforeDelete(tx *gorm.DB) (err error) {
	var list = []any{
		UserArticleFavorModel{},
		ArticleModel{},
		CategoryModel{},
		CollectModel{},
		CommentModel{},
		UserCommentFavorModel{},
		LogModel{},
		UserArticleCollectModel{},
		UserArticleHistoryModel{},
		UserChatActionModel{},
		UserFocusModel{},
		UserGlobalNotificationModel{},
		UserLoginModel{},
		UserTopArticleModel{},
	}
	for _, model := range list {
		count := tx.Delete(&model, "user_id = ?", u.ID).RowsAffected
		logrus.Infof("删除 %s 成功%d条", reflect.TypeOf(model).Name(), count)
	}

	var chatList []ChatModel
	tx.Find(&chatList, "send_user_id = ? or rev_user_id = ?", u.ID, u.ID).Delete(&chatList)
	logrus.Infof("删除关联对话 %d条", len(chatList))

	var messageList []MessageModel
	tx.Find(&messageList, "rev_user_id = ?", u.ID).Delete(&messageList)
	logrus.Infof("删除关联消息 %d条", len(messageList))

	return nil
}

// CodeAge 计算码龄
func (u *UserModel) CodeAge() int {
	sub := time.Now().Sub(u.CreatedAt)
	return int(math.Ceil(sub.Hours() / 24 / 365))
}

// UserConfModel 用户配置表
type UserConfModel struct {
	UserID         uint       `gorm:"primaryKey;unique" json:"userID"`
	UserModel      UserModel  `gorm:"foreignKey:UserID" json:"-"`
	InterestTags   []string   `gorm:"type:text;serializer:json" json:"interest_tags"` //兴趣标签
	UpdateMark     *time.Time `json:"update_mark"`                                    //上次修改用户配置时间
	OpenCollection bool       `json:"open_collection"`                                //公开收藏
	OpenFans       bool       `json:"open_fans"`                                      //公开粉丝
	OpenFollows    bool       `json:"open_follows"`                                   //公开关注
	HomeStyle      uint       `json:"home_style"`                                     //主页样式id
	ViewCount      int        `json:"views_count"`                                    // 主页的访问次数
}
