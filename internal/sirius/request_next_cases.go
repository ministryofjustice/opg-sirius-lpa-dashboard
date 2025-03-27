package sirius

import (
	"net/http"
)

func (c *Client) RequestNextCases(ctx Context) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/request-new-cases", nil)
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
