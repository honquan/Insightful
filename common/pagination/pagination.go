package pagination

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

var NoUse = &Pagination{NoUse: true}

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

var mpSort = map[Order]int{
	OrderAsc:  1,
	OrderDesc: -1,
}

type Sort struct {
	SortBy   string   `json:"sort_by" query:"sort_by"`
	SortMul  []string `json:"-"`
	Order    string   `json:"order" query:"order"`
	OrderMul []string `json:"-"`
}

func (s Sort) SortOption() bson.D {
	s.SortMul = strings.Split(s.SortBy, ",")
	s.OrderMul = strings.Split(s.Order, ",")
	sortOp := bson.D{}
	maxL := len(s.SortMul)
	if len(s.OrderMul) < maxL {
		maxL = len(s.OrderMul)
	}
	for i := 0; i < maxL; i++ {
		if strings.HasPrefix(s.SortMul[i], "field.") {
			key := strings.TrimPrefix(s.SortMul[i], "field.")
			sortOp = append(sortOp, bson.E{"sortMeta." + key, mpSort[Order(s.OrderMul[i])]})
			continue
		}
		sortOp = append(sortOp, bson.E{s.SortMul[i], mpSort[Order(s.OrderMul[i])]})
	}
	return sortOp
}

type Pagination struct {
	NoUse      bool  `json:"-"`
	Total      int64 `json:"total"`
	TotalPage  int64 `json:"total_page"`
	KeepOffset bool  `json:"-"`
	Page       int64 `json:"page" query:"page"`
	Limit      int64 `json:"limit" query:"limit"`
	Sort
}

func (p *Pagination) SetTotal(total int64) {
	p.Total = total
	p.TotalPage = total / p.Limit
	if total%p.Limit != 0 {
		p.TotalPage++
	}
}

func (p *Pagination) Correct() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 20
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.Order == "" {
		p.Order = string(OrderDesc)
	}
}

func (p *Pagination) Skip() int64 {
	return (p.Page - 1) * p.Limit
}

func (p *Pagination) FindOption() *options.FindOptions {
	p.Correct()
	op := options.Find().SetCollation(&options.Collation{Locale: "vi", CaseLevel: true})
	if p.NoUse {
		if p.KeepOffset {
			return op.SetLimit(p.Limit).SetSkip(p.Skip()).SetSort(p.SortOption())
		}
		return op.SetSort(p.SortOption())
	}
	return op.SetLimit(p.Limit).SetSkip(p.Skip()).SetSort(p.SortOption())
}

func (p *Pagination) FindOptionWithoutSkip() *options.FindOptions {
	p.Correct()
	op := options.Find().SetCollation(&options.Collation{Locale: "vi", CaseLevel: true})
	if p.NoUse {
		if p.KeepOffset {
			return op.SetLimit(p.Limit).SetSort(p.SortOption())
		}
		return op.SetSort(p.SortOption())
	}
	return op.SetLimit(p.Limit).SetSort(p.SortOption())
}

func (p *Pagination) Aggregation() bson.A {
	return bson.A{
		bson.D{{"$sort", p.SortOption()}},
		bson.D{{"$skip", (p.Page - 1) * p.Limit}},
		bson.D{{"$limit", p.Limit}},
	}
}
