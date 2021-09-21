package sirius

import (
	"net/http"
)

func (c *Client) RequestNextPaymentTask(ctx Context) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/request-new-payment-task", nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
