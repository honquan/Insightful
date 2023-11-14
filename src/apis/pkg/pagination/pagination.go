package pagination

//var NoUse = &Pagination{NoUse: true}
//
//type Order string
//
//const (
//	OrderAsc  Order = "asc"
//	OrderDesc Order = "desc"
//)
//
//var mpSort = map[Order]int{
//	OrderAsc:  1,
//	OrderDesc: -1,
//}
//
//type Sort struct {
//	SortBy   string   `json:"sort_by" query:"sort_by"`
//	SortMul  []string `json:"-"`
//	Order    string   `json:"order" query:"order"`
//	OrderMul []string `json:"-"`
//}
//
//type Pagination struct {
//	NoUse      bool  `json:"-"`
//	Total      int64 `json:"total"`
//	TotalPage  int64 `json:"total_page"`
//	KeepOffset bool  `json:"-"`
//	Page       int64 `json:"page" query:"page"`
//	Limit      int64 `json:"limit" query:"limit"`
//	Sort
//}
//
//func (p *Pagination) SetTotal(total int64) {
//	p.Total = total
//	p.TotalPage = total / p.Limit
//	if total%p.Limit != 0 {
//		p.TotalPage++
//	}
//}
//
//func (p *Pagination) Correct() {
//	if p.Page < 1 {
//		p.Page = 1
//	}
//	if p.Limit < 1 {
//		p.Limit = 20
//	}
//	if p.SortBy == "" {
//		p.SortBy = "created_at"
//	}
//	if p.Order == "" {
//		p.Order = string(OrderDesc)
//	}
//}
//
//func (p *Pagination) Skip() int64 {
//	return (p.Page - 1) * p.Limit
//}
