package middleware

import (
	"blogX_server/global"
	"blogX_server/utils/jwts"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"time"
)

type CacheOption struct {
	Prefix  CacheMiddlewarePrefix
	Time    time.Duration
	Params  []string
	NoCache func(c *gin.Context) bool
	IsUser  bool
}

type CacheMiddlewarePrefix string

const (
	CacheBannerPrefix        CacheMiddlewarePrefix = "cache_banner_"
	CacheDataPrefix          CacheMiddlewarePrefix = "cache_data_"
	CacheArticleDetailPrefix CacheMiddlewarePrefix = "cache_article_detail_"
)

func NewBannerCacheOption() CacheOption {
	return CacheOption{
		Prefix: CacheBannerPrefix,
		Time:   time.Hour,
		Params: []string{"type"},
		NoCache: func(c *gin.Context) bool {
			var referer = c.GetHeader("referer")
			if strings.Contains(referer, "admin") {
				// 后台来的，不走缓存
				return true
			}
			return false
		},
	}
}

func NewDataCacheOption() CacheOption {
	return CacheOption{
		Prefix: CacheDataPrefix,
		Time:   time.Minute,
	}
}
func NewArticleDetailCacheOption() CacheOption {
	return CacheOption{
		Prefix: CacheArticleDetailPrefix,
		Time:   time.Minute,
		IsUser: true,
	}
}

type CacheResponseWriter struct {
	gin.ResponseWriter
	Body []byte
}

func (w *CacheResponseWriter) Write(data []byte) (int, error) {
	w.Body = append(w.Body, data...)
	return w.ResponseWriter.Write(data)
}
func CacheMiddleware(option CacheOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		values := url.Values{}
		for _, key := range option.Params {
			values.Add(key, c.Query(key))
		}
		var key string
		if option.IsUser {
			var userID uint = 0
			claims, err := jwts.ParseTokenByGin(c)
			if err == nil && claims != nil {
				userID = claims.UserID
			}
			key = fmt.Sprintf("%s%d%s", option.Prefix, userID, values.Encode())
		} else {
			key = fmt.Sprintf("%s%s", option.Prefix, values.Encode())
		}

		// 请求部分
		val, err := global.Redis.Get(context.Background(), key).Result()
		fmt.Println(key, err)
		// （找到缓存 && 没有配置noCache ）|| (找到缓存 && noCache = false)
		if (err == nil && option.NoCache == nil) || (err == nil && option.NoCache(c) == false) {
			c.Abort()
			fmt.Println("走缓存了")
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.Writer.Write([]byte(val))
			return
		}
		w := &CacheResponseWriter{
			ResponseWriter: c.Writer,
		}
		c.Writer = w
		c.Next()
		// 响应
		body := string(w.Body)
		// 加入到缓存里面
		global.Redis.Set(context.Background(), key, body, option.Time)
	}
}

func CacheClose(prefix CacheMiddlewarePrefix) {
	keys, err := global.Redis.Keys(context.Background(), fmt.Sprintf("%s*", prefix)).Result()
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}
	if len(keys) > 0 {
		logrus.Infof("删除前缀 %s 缓存 共 %d 条", prefix, len(keys))
		global.Redis.Del(context.Background(), keys...)
	}

}
