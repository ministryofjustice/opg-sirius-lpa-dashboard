package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CasesByTeam struct {
	Cases      []Case
	Stats      CasesByTeamMetadata
	Pagination *Pagination
}

type CasesByTeamMetadata struct {
	WorkedTotal    int                         `json:"workedTotal"`
	Worked         []CasesByTeamMetadataMember `json:"worked"`
	TasksCompleted []CasesByTeamMetadataMember `json:"tasksCompleted"`
}

type CasesByTeamMetadataMember struct {
	Assignee Assignee `json:"assignee"`
	Total    int      `json:"total"`
}

func (c *Client) CasesByTeam(ctx Context, id int, criteria Criteria) (*CasesByTeam, error) {
	url := fmt.Sprintf("/api/v1/teams/%d/cases?%s", id, criteria.String())

	req, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v struct {
		Pages    apiPages            `json:"pages"`
		Total    int                 `json:"total"`
		Limit    int                 `json:"limit"`
		Cases    []Case              `json:"cases"`
		Metadata CasesByTeamMetadata `json:"metadata"`
	}

	err = json.NewDecoder(resp.Body).Decode(&v)

	return &CasesByTeam{
		Cases: v.Cases,
		Stats: v.Metadata,
		Pagination: &Pagination{
			TotalItems:  v.Total,
			CurrentPage: v.Pages.Current,
			TotalPages:  v.Pages.Total,
			PageSize:    v.Limit,
		},
	}, err
}
