package slicex

// Exist 判断切片中是否存在元素
func Exist[T comparable](s []T, e T) bool {
	for _, i := range s {
		if i == e {
			return true
		}
	}
	return false
}
