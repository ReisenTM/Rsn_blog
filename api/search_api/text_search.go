package search_api

import (
	"blogX_server/common"
	"blogX_server/common/resp"
	"blogX_server/global"
	"blogX_server/middleware"
	"blogX_server/model"
	"blogX_server/service/text_service"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type TextSearchRequest struct {
	common.PageInfo
}
type TextSearchResponse struct {
	ArticleID uint   `json:"article_id"`
	Head      string `json:"head"`
	Body      string `json:"body"`
	Flag      string `json:"flag"`
}

func (SearchApi) TextSearchView(c *gin.Context) {
	cr := middleware.GetBind[TextSearchRequest](c)
	if global.EsClient == nil {
		//如果没配置es，降级
		_list, count, _ := common.ListQuery(model.TextModel{}, common.Options{
			PageInfo: cr.PageInfo,
			Likes:    []string{"head", "body"},
		})

		var list = make([]TextSearchResponse, 0)
		for _, mod := range _list {
			list = append(list, TextSearchResponse{
				ArticleID: mod.ArticleID,
				Head:      mod.Head,
				Body:      mod.Body,
				Flag:      mod.Head,
			})
		}
		resp.OkWithList(list, count, c)
		return
	}

	query := elastic.NewBoolQuery()
	if cr.Key != "" {
		query.Should(
			elastic.NewMatchQuery("head", cr.Key),
			elastic.NewMatchQuery("body", cr.Key),
		)
	}
	highlight := elastic.NewHighlight()
	highlight.Field("head")
	highlight.Field("body")
	result, err := global.EsClient.
		Search(model.TextModel{}.Index()).
		Query(query).
		Highlight(highlight).
		From(cr.GetOffset()).
		Size(cr.GetLimit()).
		Do(context.Background())
	if err != nil {
		source, _ := query.Source()
		byteData, _ := json.Marshal(source)
		logrus.Errorf("查询失败 %s \n %s", err, string(byteData))
		resp.FailWithMsg("查询失败", c)
		return
	}
	count := result.Hits.TotalHits.Value
	var list = make([]TextSearchResponse, 0)
	for _, hit := range result.Hits.Hits {
		var item text_service.TextModel
		err = json.Unmarshal(hit.Source, &item)
		if err != nil {
			logrus.Warnf("解析失败 %s  %s", err, string(hit.Source))
			continue
		}
		head := item.Head
		if len(hit.Highlight["head"]) > 0 {
			//默认高亮第一个找到的元素
			item.Head = hit.Highlight["head"][0]
		}
		if len(hit.Highlight["body"]) > 0 {
			item.Body = hit.Highlight["body"][0]
		}
		list = append(list, TextSearchResponse{
			ArticleID: item.ArticleID,
			Head:      item.Head,
			Body:      item.Body,
			Flag:      head,
		})
	}
	resp.OkWithList(list, int(count), c)
}
