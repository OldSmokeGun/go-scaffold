package slicex

// InStringSlice 判断字符串是否在切片中
func InStringSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
