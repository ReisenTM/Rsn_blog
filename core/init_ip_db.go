package core

import (
	ipUtils "blogX_server/utils/ip"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
	"strings"
)

var searcher *xdb.Searcher

const LOCFOMMAT = 5

func InitIPDB() {
	var dbPath = "init/ip2region.xdb"
	_searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		logrus.Fatalf("ip地址数据库加载失败: %s\n", err)
		return
	}
	//不关闭因为后面还需要用
	//defer searcher.Close()
	searcher = _searcher
}

func GetIPLoc(ip string) (location string) {
	//利用区间先快速判断是否是内网
	if ipUtils.HasLocalIPAddr(ip) {
		return "内网"
	}
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		logrus.Warnf("错误的ip地址:[%s]", ip)
		return "异常地址"
	}
	//处理addrList
	_addrList := strings.Split(region, "|")
	if len(_addrList) != LOCFOMMAT {
		//出现概率目前极低
		logrus.Warnf("异常的ip地址:[%s]", ip)
		return "未知地址"
	}
	//_addrList五个部分
	//国家|0|省|市｜运营商
	country := _addrList[0]
	province := _addrList[2]
	city := _addrList[3]
	//为了实现打印最大精度
	//例如：如果最大精度为省，那么只打印 中国.河南
	if province != "0" && city != "0" {
		return fmt.Sprintf("%s·%s", country, city)
	}
	if country != "0" && province != "0" {
		return fmt.Sprintf("%s·%s", country, province)
	}
	if country != "0" {
		return country
	}
	return region
}
