package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const PotUserEmail string =  "opgcasework@publicguardian.gov.uk"

func (c *Client) User(ctx Context, id int) (Assignee, error) {
	var v Assignee

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/users/%d", id), nil)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
