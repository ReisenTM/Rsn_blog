package data_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/model/enum"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type GrowthDataRequest struct {
	Type int8 `form:"type" binding:"required,oneof=1 2 3"` //1 流量 2 文章 3 用户
}

type GrowthDataResponse struct {
	GrowthRate int      `json:"growth_rate"` // 增长率
	GrowthNum  int      `json:"growth_num"`  // 增长数
	CountList  []int    `json:"count_list"`
	DateList   []string `json:"date_list"`
}
type Table struct {
	Date  string `gorm:"column:date"`
	Count int    `gorm:"column:count"`
}

// GrowthDataView 七天增长数据
func (DataApi) GrowthDataView(c *gin.Context) {
	cr := middleware.GetBind[GrowthDataRequest](c)

	now := time.Now()
	before7 := now.AddDate(0, 0, -7)
	var dataList []Table

	switch cr.Type {
	case 1:
		//流量增长
		global.DB.Model(model.SiteFlowModel{}).Where("created_at >= ? and created_at <= ?",
			before7.Format("2006-01-02")+" 00:00:00",
			now.Format("2006-01-02 15:04:05"),
		).
			Select("date(created_at) as date", "sum(count) as count").
			Group("date").Scan(&dataList)
	case 2:
		//发布文章增长
		global.DB.Model(model.ArticleModel{}).Where("created_at >= ? and created_at <= ? and status = ?",
			before7.Format("2006-01-02")+" 00:00:00",
			now.Format("2006-01-02 15:04:05"),
			enum.ArticleStatusPublished).
			Select("date(created_at) as date", "count(id) as count").
			Group("date").Scan(&dataList)
	case 3:
		//用户量增长
		global.DB.Model(model.UserModel{}).Where("created_at >= ? and created_at <= ?",
			before7.Format("2006-01-02")+" 00:00:00",
			now.Format("2006-01-02 15:04:05"),
		).
			Select("date(created_at) as date", "count(id) as count").
			Group("date").Scan(&dataList)
	}
	var dateMap = map[string]int{}
	for _, model := range dataList {
		date := strings.Split(model.Date, "T")[0]
		dateMap[date] = model.Count
	}

	response := GrowthDataResponse{}
	for i := 0; i < 7; i++ {
		date := before7.AddDate(0, 0, i+1)
		dateS := date.Format("2006-01-02")
		count, _ := dateMap[dateS]
		response.CountList = append(response.CountList, count)
		response.DateList = append(response.DateList, dateS)
	}
	// 算增长，找最后一个和最后一个的前一个
	response.GrowthNum = response.CountList[6] - response.CountList[5]
	if response.CountList[5] == 0 {
		response.GrowthRate = 100
	} else {
		response.GrowthRate = int(float64(response.GrowthNum) / float64(response.CountList[5]) * 100)
	}
	resp.OkWithData(response, c)
}
