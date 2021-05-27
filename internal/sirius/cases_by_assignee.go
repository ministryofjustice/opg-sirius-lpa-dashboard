package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Case struct {
	ID          int        `json:"id"`
	Uid         string     `json:"uId"`
	Donor       Donor      `json:"donor"`
	SubType     string     `json:"caseSubtype"`
	ReceiptDate SiriusDate `json:"receiptDate"`
	Status      string     `json:"status"`
	TaskCount   int        `json:"taskCount"`
}

type Donor struct {
	ID        int    `json:"id"`
	Uid       string `json:"uId"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
}

func (d Donor) DisplayName() string {
	return d.Firstname + " " + d.Surname
}

type SortOrder string

const (
	Ascending  SortOrder = "asc"
	Descending SortOrder = "desc"
)

type CasesByAssigneeFilter struct {
	Status string
}

type CasesByAssigneeSort struct {
	Field string
	Order SortOrder
}

type CasesByAssigneeCriteria struct {
	Page   int
	Limit  int
	Filter CasesByAssigneeFilter
	Sort   CasesByAssigneeSort
}

func (c *Client) CasesByAssignee(ctx Context, id int, criteria CasesByAssigneeCriteria) ([]Case, *Pagination, error) {
	filter := "caseType:lpa,active:true"
	if criteria.Filter.Status != "" {
		filter = fmt.Sprintf("%s,status:%s", filter, criteria.Filter.Status)
	}

	sortField := "receiptDate"
	if criteria.Sort.Field != "" {
		sortField = criteria.Sort.Field
	}

	sortOrder := Ascending
	if criteria.Sort.Order != "" {
		sortOrder = criteria.Sort.Order
	}

	page := 1
	if criteria.Page != 0 {
		page = criteria.Page
	}

	limit := 25
	if criteria.Limit != 0 {
		limit = criteria.Limit
	}

	url := fmt.Sprintf("/api/v1/assignees/%d/cases?page=%d&filter=%s&limit=%d&sort=%s:%s", id, page, filter, limit, sortField, sortOrder)

	req, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, newStatusError(resp)
	}

	var v struct {
		Pages apiPages `json:"pages"`
		Total int      `json:"total"`
		Limit int      `json:"limit"`
		Cases []Case   `json:"cases"`
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v.Cases, &Pagination{
		TotalItems:  v.Total,
		CurrentPage: v.Pages.Current,
		TotalPages:  v.Pages.Total,
		PageSize:    v.Limit,
	}, err
}

func (c *Client) HasWorkableCase(ctx Context, id int) (bool, error) {
	_, pagination, err := c.CasesByAssignee(ctx, id, CasesByAssigneeCriteria{
		Filter: CasesByAssigneeFilter{
			Status: "Pending",
		},
		Page: 1,
	})

	return pagination.TotalItems > 0, err
}
