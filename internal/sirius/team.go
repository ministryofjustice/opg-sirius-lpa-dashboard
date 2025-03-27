package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) Team(ctx Context, id int) (Team, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/teams/%d", id), nil)
	if err != nil {
		return Team{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return Team{}, err
	}
	defer resp.Body.Close() //nolint:errcheck // no need to check error when closing body

	if resp.StatusCode == http.StatusUnauthorized {
		return Team{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return Team{}, newStatusError(resp)
	}

	var v apiTeam
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return Team{}, err
	}

	team := Team{
		ID:          v.ID,
		DisplayName: v.DisplayName,
	}

	for _, m := range v.Members {
		team.Members = append(team.Members, TeamMember{
			ID:          m.ID,
			DisplayName: m.DisplayName,
		})
	}

	return team, nil
}
