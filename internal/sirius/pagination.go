package sirius

type Pagination struct {
	TotalItems  int
	CurrentPage int
	TotalPages  int
	PageSize    int
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
