package request

type PageSearch struct {
	PageNum  int64  `json:"pageNum,omitempty" form:"pageNum"`   // 页码
	PageSize int64  `json:"pageSize,omitempty" form:"pageSize"` // 每页显示数量
	Keyword  string `json:"keyword,omitempty" form:"keyword"`   // 关键词
}

type Ids struct {
	Ids []int `json:"ids,omitempty" form:"ids"` // id切片
}
