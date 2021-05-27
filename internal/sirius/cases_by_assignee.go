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

func (c *Client) CasesByAssignee(ctx Context, id int, criteria Criteria) ([]Case, *Pagination, error) {
	criteria = criteria.Filter("caseType", "lpa").Filter("active", "true")

	url := fmt.Sprintf("/api/v1/assignees/%d/cases?%s", id, criteria.String())

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
	_, pagination, err := c.CasesByAssignee(ctx, id, Criteria{}.Filter("status", "Pending").Page(1))

	return pagination.TotalItems > 0, err
}
