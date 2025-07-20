package data_api

import (
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/model"
	"blogX_server/model/enum"
	"github.com/gin-gonic/gin"
	"time"
)

type ArticleYearDataResponse struct {
	GrowthRate int      `json:"growth_rate"` // 增长率
	GrowthNum  int      `json:"growth_num"`  // 增长数
	CountList  []int    `json:"count_list"`
	DateList   []string `json:"date_list"`
}

func (DataApi) ArticleYearDataView(c *gin.Context) {
	// 1 2 3 4 5 6 7
	// 1 10%
	now := time.Now()
	// 12月前的时间
	before12Month := now.AddDate(0, -12, 0)
	var dataList []Table

	// 查询七天内的文章
	global.DB.Model(model.ArticleModel{}).Where("created_at >= ? and created_at <= ? and status = ?",
		before12Month.Format("2006-01-02")+" 00:00:00",
		now.Format("2006-01-02 15:04:05"),
		enum.ArticleStatusPublished).
		Select("month(created_at) as date", "count(id) as count").
		Group("date").Scan(&dataList) //按照月份分组统计文章数量

	var dateMap = map[string]int{}
	for _, model := range dataList {
		date := model.Date
		dateMap[date] = model.Count
	}

	response := ArticleYearDataResponse{}
	for i := 0; i < 12; i++ {
		date := before12Month.AddDate(0, i+1, 0)
		dateS := date.Format("1")
		count, _ := dateMap[dateS]
		response.CountList = append(response.CountList, count)
		response.DateList = append(response.DateList, date.Format("2006-01"))
	}
	// 算增长，找最后一个和最后一个的前一个
	response.GrowthNum = response.CountList[11] - response.CountList[10]
	if response.CountList[10] == 0 {
		response.GrowthRate = 100
	} else {
		response.GrowthRate = int(float64(response.GrowthNum) / float64(response.CountList[10]) * 100)
	}
	resp.OkWithData(response, c)
}
