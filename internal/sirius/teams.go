package sirius

import (
	"encoding/json"
	"net/http"
)

type apiTeam struct {
	ID          int       `json:"id"`
	DisplayName string    `json:"displayName"`
	TeamType    *struct{} `json:"teamType"`
	Members     []struct {
		ID          int    `json:"id"`
		DisplayName string `json:"displayName"`
	} `json:"members"`
}

type Team struct {
	ID          int
	DisplayName string
	Members     []TeamMember
}

type TeamMember struct {
	ID          int
	DisplayName string
}

func (c *Client) Teams(ctx Context) ([]Team, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/teams", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //#nosec G307 false positive

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v []apiTeam
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	var teams []Team

	for _, t := range v {
		if t.TeamType != nil {
			continue
		}

		team := Team{
			ID:          t.ID,
			DisplayName: t.DisplayName,
			Members:     make([]TeamMember, len(t.Members)),
		}

		for i, m := range t.Members {
			team.Members[i] = TeamMember{
				ID:          m.ID,
				DisplayName: m.DisplayName,
			}
		}

		teams = append(teams, team)
	}

	return teams, nil
}
