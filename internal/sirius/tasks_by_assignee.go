package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Task struct {
	ID        int            `json:"id"`
	Status    string         `json:"status"`
	DueDate   SiriusDate     `json:"dueDate"`
	Name      string         `json:"name"`
	CaseItems []TaskCaseItem `json:"caseItems"`
}

func (t Task) Case() TaskCaseItem {
	return t.CaseItems[0]
}

type TaskCaseItem struct {
	ID      int    `json:"id"`
	Uid     string `json:"uid"`
	Donor   Donor  `json:"donor"`
	SubType string `json:"caseSubtype"`
}

func (c *Client) TasksByAssignee(ctx Context, id int, criteria Criteria) ([]Task, *Pagination, error) {
	url := fmt.Sprintf("/api/v1/assignees/%d/tasks?%s", id, criteria.String())

	req, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

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
		Tasks []Task   `json:"tasks"`
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v.Tasks, &Pagination{
		TotalItems:  v.Total,
		CurrentPage: v.Pages.Current,
		TotalPages:  v.Pages.Total,
		PageSize:    v.Limit,
	}, err
}
