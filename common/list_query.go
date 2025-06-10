package common

import (
	"blogX_server/global"
	"fmt"
	"gorm.io/gorm"
)

type PageInfo struct {
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
	Key   string `form:"key"`
	Order string `form:"order"` // 前端可以覆盖
}

// GetPage 获取页数
func (p PageInfo) GetPage() int {
	if p.Page > 20 || p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetLimit 每页限制显示数量
func (p PageInfo) GetLimit() int {
	if p.Limit <= 0 || p.Limit > 100 {
		return 10
	}
	return p.Limit
}

// GetOffset 获取页数偏移量
func (p PageInfo) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

type Options struct {
	PageInfo     PageInfo
	Likes        []string
	Preloads     []string
	Where        *gorm.DB //自定义查询：如复杂sql语句
	Debug        bool
	DefaultOrder string
}

// ListQuery 列表查询
func ListQuery[T any](model T, option Options) (list []T, count int, err error) {
	//基础查询
	query := global.DB.Model(model).Where(model)
	//日志
	if option.Debug {
		query = query.Debug()
	}
	//模糊匹配
	if len(option.Likes) > 0 && option.PageInfo.Key != "" {
		likes := global.DB.Where("")
		for _, like := range option.Likes {
			likes = likes.Or(fmt.Sprintf("%s LIKE ?", like),
				"%"+option.PageInfo.Key+"%")
		}
		query = query.Where(likes)
	}
	//高级查询
	if option.Where != nil {
		query = query.Where(option.Where)
	}
	//预加载
	for _, preload := range option.Preloads {
		query = query.Preload(preload)
	}
	//统计数量
	var _count int64
	query.Count(&_count)
	count = int(_count)
	//分页
	limit := option.PageInfo.GetLimit()
	offset := option.PageInfo.GetOffset()
	query = query.Limit(limit).Offset(offset)
	//排序
	if option.PageInfo.Order != "" {
		//在外层配置了
		query = query.Order(option.PageInfo.Order)
	} else {
		//否则按默认排序
		if option.DefaultOrder != "" {
			query = query.Order(option.DefaultOrder)
		}
	}

	err = query.Find(&list).Error
	return
}
