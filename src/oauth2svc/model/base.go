package model

type OrderBy string

var (
	DESC OrderBy = "desc"
	ASC  OrderBy = "asc"
)

func (ob OrderBy) String() string {
	if ob == DESC ||
		ob == ASC {
		return string(ob)
	}
	return string(DESC)
}

func (ob OrderBy) Bool() bool {
	return ob == DESC
}

type Pagination struct {
	Limit   int     `form:"limit" validate:"required,min=1,max=250"`
	Offset  int     `form:"offset"`
	SortBy  string  `form:"sort_by"`
	OrderBy OrderBy `form:"order_by"`
}

type PageInfo struct {
	Limit     int `form:"limit" json:"limit"`
	Offset    int `form:"offset" json:"offset"`
	TotalData int `form:"total_data" json:"total_data"`
}
