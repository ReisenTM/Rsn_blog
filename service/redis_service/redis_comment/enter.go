package redis_comment

import (
	"blogX_server/global"
	"context"
	"github.com/sirupsen/logrus"
	"strconv"
)

type commentCacheType string

const (
	commentCacheReply commentCacheType = "comment_reply_key"
	commentCacheFavor commentCacheType = "comment_favor_key"
)

func set(t commentCacheType, commentID uint, n int) {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(commentID))).Int()
	num += n
	global.Redis.HSet(context.Background(), string(t), strconv.Itoa(int(commentID)), num)
}
func SetCacheReply(commentID uint, n int) {
	set(commentCacheReply, commentID, n)
}
func SetCacheFavor(commentID uint, n int) {
	set(commentCacheFavor, commentID, n)
}
func get(t commentCacheType, commentID uint) int {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(commentID))).Int()
	return num
}
func GetCacheReply(commentID uint) int {
	return get(commentCacheReply, commentID)
}
func GetCacheFavor(commentID uint) int {
	return get(commentCacheFavor, commentID)
}
func GetAll(t commentCacheType) (mp map[uint]int) {
	res, err := global.Redis.HGetAll(context.Background(), string(t)).Result()
	if err != nil {
		return
	}
	mp = make(map[uint]int)
	for k, v := range res {
		num, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		k, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		mp[uint(k)] = num
	}
	return mp
}
func GetAllCacheReply() (mps map[uint]int) {
	return GetAll(commentCacheReply)
}
func GetAllCacheFavor() (mps map[uint]int) {
	return GetAll(commentCacheFavor)
}

func Clear() {
	err := global.Redis.Del(context.Background(), "comment_reply_key", "comment_favor_key").Err()
	if err != nil {
		logrus.Error(err)
	}
}
