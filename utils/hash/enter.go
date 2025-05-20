package hash

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// Md5 字符串hash
func Md5(data []byte) string {
	md5New := md5.New()
	md5New.Write(data)
	return hex.EncodeToString(md5New.Sum(nil))
}

// FileHash 文件hash
func FileHash(filename string) (hash string, err error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	hash = Md5(bytes)
	return hash, nil
}
