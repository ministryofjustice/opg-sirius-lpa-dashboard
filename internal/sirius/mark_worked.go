package sirius

import (
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) MarkWorked(ctx Context, id int) error {
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/lpas/%d", id), strings.NewReader(`{"worked":true}`))
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
