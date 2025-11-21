package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type feedbackRequest struct {
	Message string `json:"message"`
}

func (c *Client) Feedback(ctx Context, message string) error {
	data := feedbackRequest{
		Message: message,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/lpa-api/v1/feedback/poas", &buf)
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
		var v struct {
			Detail string `json:"detail"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ClientError(v.Detail)
		}

		return newStatusError(resp)
	}

	return nil
}
