package server

import "github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"

type Pagination struct {
	Query       string
	TotalItems  int
	CurrentPage int
	TotalPages  int
	PageSize    int
}

func newPagination(p *sirius.Pagination) *Pagination {
	return newPaginationWithQuery(p, "")
}

func newPaginationWithQuery(p *sirius.Pagination, q string) *Pagination {
	if p == nil {
		return nil
	}

	if q == "" {
		q = "?"
	} else {
		q = "?" + q + "&"
	}

	return &Pagination{
		Query:       q,
		TotalItems:  p.TotalItems,
		CurrentPage: p.CurrentPage,
		TotalPages:  p.TotalPages,
		PageSize:    p.PageSize,
	}
}

func (p *Pagination) Start() int {
	return (p.CurrentPage-1)*p.PageSize + 1
}

func (p *Pagination) End() int {
	end := p.CurrentPage * p.PageSize
	if end < p.TotalItems {
		return end
	}

	return p.TotalItems
}

func (p *Pagination) HasPrevious() bool {
	return p.CurrentPage > 1
}

func (p *Pagination) PreviousPage() int {
	return p.CurrentPage - 1
}

func (p *Pagination) HasNext() bool {
	return p.TotalItems > p.CurrentPage*p.PageSize
}

func (p *Pagination) NextPage() int {
	return p.CurrentPage + 1
}

func (p *Pagination) Pages() []int {
	pages := make([]int, p.TotalPages)
	for i := 0; i < p.TotalPages; i++ {
		pages[i] = i + 1
	}
	return pages
}
