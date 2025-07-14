package search_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"sort"
)

type TagAggResponse struct {
	Tag          string `json:"tag"`
	ArticleCount int    `json:"article_count"`
}

// TagAggView 标签关联文章的聚合
func (SearchApi) TagAggView(c *gin.Context) {
	var cr = middleware.GetBind[common.PageInfo](c)
	var list = make([]TagAggResponse, 0)
	if global.EsClient == nil {
		var articleList []model.ArticleModel
		global.DB.Find(&articleList, "tags <> ''")
		var tagMap = map[string]int{}
		//数量统计
		for _, model := range articleList {
			for _, tag := range model.Tags {
				count, ok := tagMap[tag]
				if !ok {
					tagMap[tag] = 1
					continue
				}
				tagMap[tag] = count + 1
			}
		}
		for tag, count := range tagMap {
			list = append(list, TagAggResponse{
				Tag:          tag,
				ArticleCount: count,
			})
		}
		sort.Slice(list, func(i, j int) bool {
			//按降序排列
			return list[i].ArticleCount > list[j].ArticleCount
		})
		resp.OkWithList(list, len(list), c)
		return
	}

	agg := elastic.NewTermsAggregation().Field("tags")
	agg.SubAggregation("page",
		elastic.NewBucketSortAggregation().
			From(cr.GetOffset()).
			Size(cr.Limit))
	query := elastic.NewBoolQuery()
	query.MustNot(elastic.NewTermQuery("tags", ""))
	result, err := global.EsClient.
		Search(model.ArticleModel{}.Index()).
		Query(query).
		Aggregation("tags", agg).
		Aggregation("tags1", elastic.NewCardinalityAggregation().Field("tags")).
		Size(0).Do(context.Background())
	if err != nil {
		logrus.Errorf("查询失败 %s", err)
		resp.FailWithMsg("查询失败", c)
		return
	}
	var t AggType
	var val = result.Aggregations["tags"]
	err = json.Unmarshal(val, &t)
	if err != nil {
		logrus.Errorf("解析json失败 %s %s", err, string(val))
		resp.FailWithMsg("查询失败", c)
		return
	}
	var co Agg1Type
	err = json.Unmarshal(result.Aggregations["tags1"], &co)
	if err != nil {
		logrus.Errorf("Json Unmarshal %v", err)
		resp.FailWithMsg("查询失败", c)
		return
	}
	for _, bucket := range t.Buckets {
		list = append(list, TagAggResponse{
			Tag:          bucket.Key,
			ArticleCount: bucket.DocCount,
		})
	}
	resp.OkWithList(list, co.Value, c)
	return
}

type AggType struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
	} `json:"buckets"`
}
type Agg1Type struct {
	Value int `json:"value"`
}
