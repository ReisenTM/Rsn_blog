package redis_site

import (
	"blogX_server/global"
	"context"
)

const key = "blogx_site_flow"

//网站流量统计

func SetFlow() {
	v, _ := global.Redis.Get(context.Background(), key).Int()
	global.Redis.Set(context.Background(), key, v+1, 0)
}

func GetFlow() int {
	v, _ := global.Redis.Get(context.Background(), key).Int()
	return v
}

func ClearFlow() {
	global.Redis.Del(context.Background(), key)
}
