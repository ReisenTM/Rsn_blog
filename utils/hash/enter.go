package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(data []byte) string {
	md5New := md5.New()
	md5New.Write(data)
	return hex.EncodeToString(md5New.Sum(nil))
}
