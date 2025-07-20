package redis_user

import (
	"blogX_server/global"
	"context"
	"github.com/sirupsen/logrus"
	"strconv"
)

type userCacheType string

const (
	userCacheLook userCacheType = "user_look_key"
)

func set(t userCacheType, userID uint, n int) {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(userID))).Int()
	num += n
	global.Redis.HSet(context.Background(), string(t), strconv.Itoa(int(userID)), num)
}

func SetCacheLook(userID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(userCacheLook, userID, n)
}

func get(t userCacheType, userID uint) int {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(userID))).Int()
	return num
}
func GetCacheLook(userID uint) int {
	return get(userCacheLook, userID)
}

func GetAll(t userCacheType) (mps map[uint]int) {
	res, err := global.Redis.HGetAll(context.Background(), string(t)).Result()
	if err != nil {
		return
	}
	mps = make(map[uint]int)
	for key, numS := range res {
		iK, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		iN, err := strconv.Atoi(numS)
		if err != nil {
			continue
		}
		mps[uint(iK)] = iN
	}

	return mps
}

func GetAllCacheLook() (mps map[uint]int) {
	return GetAll(userCacheLook)
}

func Clear() {
	err := global.Redis.Del(context.Background(), "user_look_key").Err()
	if err != nil {
		logrus.Error(err)
	}
}
