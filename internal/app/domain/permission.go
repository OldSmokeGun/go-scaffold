package domain

type Permission struct {
	ID       int64  `json:"id"`
	Key      string `json:"key"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	ParentID int64  `json:"parentID"`
}
