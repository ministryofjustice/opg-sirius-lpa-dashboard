package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID int
}

type authUserResponse struct {
	ID int `json:"id"`
}

func (c *Client) UserByEmail(ctx Context, email string) (User, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/users?email=%s", email), nil)
	if err != nil {
		return User{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

	if resp.StatusCode == http.StatusUnauthorized {
		return User{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return User{}, newStatusError(resp)
	}

	var v authUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return User{}, err
	}

	return User(v), err
}
