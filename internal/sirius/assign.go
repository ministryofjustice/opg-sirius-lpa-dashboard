package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type assignRequest struct {
	Data []assignRequestItem `json:"data"`
}

type assignRequestItem struct {
	AssigneeID int    `json:"assigneeId"`
	CaseType   string `json:"caseType"`
	ID         int    `json:"id"`
}

func (c *Client) Assign(ctx Context, cases []int, assignee int) error {
	var data assignRequest
	caseList := make([]string, len(cases))

	for i, c := range cases {
		caseList[i] = strconv.Itoa(c)

		data.Data = append(data.Data, assignRequestItem{
			AssigneeID: assignee,
			CaseType:   "LPA",
			ID:         c,
		})
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return err
	}

	url := fmt.Sprintf("/api/v1/users/%d/cases/%s", assignee, strings.Join(caseList, "+"))

	req, err := c.newRequest(ctx, http.MethodPut, url, &buf)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
