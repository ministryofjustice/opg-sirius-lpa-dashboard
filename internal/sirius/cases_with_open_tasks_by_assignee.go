package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiPages struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

func (c *Client) CasesWithOpenTasksByAssignee(ctx Context, id, page int) ([]Case, *Pagination, error) {
	url := fmt.Sprintf("/api/v1/assignees/%d/cases-with-open-tasks?page=%d", id, page)

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
