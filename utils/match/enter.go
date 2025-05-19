package match

// InList 查询文件是否在列表里
func InList[T comparable](key T, list []T) bool {
	for _, s := range list {
		if key == s {
			return true
		}
	}
	return false
}
