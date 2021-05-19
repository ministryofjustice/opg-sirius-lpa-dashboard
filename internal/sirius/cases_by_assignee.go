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
}

type Donor struct {
	ID        int    `json:"id"`
	Uid       string `json:"uId"`
	FirstName string `json:"firstname"`
	Surname   string `json:"surname"`
}

func (d Donor) DisplayName() string {
	return d.FirstName + " " + d.Surname
}

func (c *Client) CasesByAssignee(ctx Context, id int) ([]Case, error) {
	var v struct {
		Cases []Case `json:"cases"`
	}

	url := fmt.Sprintf("/api/v1/assignees/%d/cases?filter=caseType:lpa,status:Pending,active:true&sort=caseSubType:asc", id)

	req, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return v.Cases, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v.Cases, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v.Cases, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v.Cases, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v.Cases, err
}
