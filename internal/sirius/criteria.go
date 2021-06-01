package sirius

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type sortOrder string

const (
	Ascending  sortOrder = "asc"
	Descending sortOrder = "desc"
)

type filter struct {
	field string
	value string
}

func (f filter) String() string {
	return fmt.Sprintf("%s:%s", f.field, f.value)
}

type sort struct {
	field string
	order sortOrder
}

func (f sort) String() string {
	return fmt.Sprintf("%s:%s", f.field, f.order)
}

type Criteria struct {
	page   int
	limit  int
	filter []filter
	sort   []sort
}

func (c Criteria) Page(id int) Criteria {
	c.page = id
	return c
}

func (c Criteria) Limit(id int) Criteria {
	c.limit = id
	return c
}

func (c Criteria) Filter(field string, value string) Criteria {
	c.filter = append(c.filter, filter{
		field: field,
		value: value,
	})
	return c
}

func (c Criteria) Sort(field string, order sortOrder) Criteria {
	c.sort = append(c.sort, sort{
		field: field,
		order: order,
	})
	return c
}

func (c *Criteria) String() string {
	params := url.Values{}

	if len(c.filter) > 0 {
		var filters []string
		for _, filter := range c.filter {
			filters = append(filters, filter.String())
		}
		params.Add("filter", strings.Join(filters, ","))
	}

	if len(c.sort) > 0 {
		var sorts []string
		for _, sort := range c.sort {
			sorts = append(sorts, sort.String())
		}
		params.Add("sort", strings.Join(sorts, ","))
	}

	if c.page != 0 {
		params.Add("page", strconv.Itoa(c.page))
	}

	if c.limit != 0 {
		params.Add("limit", strconv.Itoa(c.limit))
	}

	return params.Encode()
}
