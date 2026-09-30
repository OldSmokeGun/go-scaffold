package pagenation

// Unlimit 表示不限制条数。
const Unlimit = -1

// Param 分页参数。
type Param struct {
	Page  int `json:"page" query:"page" form:"page"`
	Limit int `json:"limit" query:"limit" form:"limit"`
}

// Result 分页结果。
type Result[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}

// Offset 计算偏移。page 从 1 开始。
func Offset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		return 0
	}

	return (page - 1) * limit
}

// Offset 返回当前参数的偏移。
func (p Param) Offset() int {
	return Offset(p.Page, p.Limit)
}

// Apply 把分页套到查询上。Limit 为 Unlimit 或小于 1 时不截断。
func Apply[Q interface {
	Offset(int) Q
	Limit(int) Q
}](q Q, p Param) Q {
	if p.Limit == Unlimit {
		return q
	}

	if p.Limit < 1 {
		return q
	}

	return q.Limit(p.Limit).Offset(Offset(p.Page, p.Limit))
}
