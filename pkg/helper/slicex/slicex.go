package slicex

// StringExist 判断切片中是否存在元素
func StringExist(s []string, e string) bool {
	for _, i := range s {
		if i == e {
			return true
		}
	}
	return false
}
