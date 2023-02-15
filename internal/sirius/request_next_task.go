package sirius

import (
	"net/http"
)

func (c *Client) RequestNextTask(ctx Context) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/request-new-task", nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
