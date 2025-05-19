package file

import (
	"blogX_server/global"
	"blogX_server/utils/match"
	"errors"
	"strings"
)

// ImageSuffixJudge 文件格式校验
func ImageSuffixJudge(filename string) (suffix string, err error) {
	_str := strings.Split(filename, ".")
	if len(_str) <= 1 {
		err = errors.New("错误的文件名")
		return
	}
	suffix = _str[len(_str)-1]
	if !match.InList(suffix, global.Config.Upload.WhiteList) {
		err = errors.New("文件非法")
		return
	}
	return
}
