package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Case struct {
	ID          int        `json:"id"`
	Uid         string     `json:"uId"`
	Donor       Donor      `json:"donor"`
	SubType     string     `json:"caseSubtype"`
	ReceiptDate SiriusDate `json:"receiptDate"`
	Status      string     `json:"status"`
	TaskCount   int        `json:"taskCount"`
	WorkedDate  SiriusDate `json:"workedDate"`
}

func (c Case) IsWorked() bool {
	day := 24 * 60 * time.Minute

	return c.WorkedDate.Truncate(day).Equal(time.Now().Truncate(day))
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

func (c *Client) CasesByAssignee(ctx Context, id int, status string, page int) ([]Case, *Pagination, error) {
	filter := "caseType:lpa,active:true"

	if status != "" {
		filter = fmt.Sprintf("%s,status:%s", filter, status)
	}

	url := fmt.Sprintf("/api/v1/assignees/%d/cases?page=%d&filter=%s&sort=caseSubType:asc", id, page, filter)

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
	_, pagination, err := c.CasesByAssignee(ctx, id, "Pending,worked:true", 1)

	return pagination.TotalItems > 0, err
}
