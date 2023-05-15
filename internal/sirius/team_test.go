package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTeam(t *testing.T) {
	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		id               int
		name             string
		setup            func()
		expectedResponse Team
		expectedError    error
	}{
		{
			name: "OK",
			id:   66,
			setup: func() {
				pact.
					AddInteraction().
					Given("LPA team with members exists").
					UponReceiving("A request for an LPA team with ID 66").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/66"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"id":          dsl.Like(66),
							"displayName": dsl.Like("Cool Team"),
							"members": dsl.EachLike(map[string]interface{}{
								"id":          dsl.Like(400),
								"displayName": dsl.Like("Carline"),
							}, 1),
						}),
					})
			},
			expectedResponse: Team{
				ID:          66,
				DisplayName: "Cool Team",
				Members: []TeamMember{
					{
						ID:          400,
						DisplayName: "Carline",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				team, err := client.Team(Context{Context: context.Background()}, tc.id)
				assert.Equal(t, tc.expectedResponse, team)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestTeamIgnored(t *testing.T) {
	pact := newIgnoredPact()
	defer pact.Teardown()

	testCases := []struct {
		id               int
		name             string
		setup            func()
		expectedResponse Team
		expectedError    error
	}{
		{
			name: "OK",
			id:   67,
			setup: func() {
				pact.
					AddInteraction().
					Given("LPA team with members exists").
					UponReceiving("A request for an LPA team with ID 67").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/67"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"id":          dsl.Like(67),
							"displayName": dsl.Like("Nottingham casework team"),
							"members": dsl.EachLike(map[string]interface{}{
								"id":          dsl.Like(600),
								"displayName": dsl.Like("Jet"),
							}, 1),
						}),
					})
			},
			expectedResponse: Team{
				ID:          67,
				DisplayName: "Nottingham casework team",
				Members: []TeamMember{
					{
						ID:          600,
						DisplayName: "Jet",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				team, err := client.Team(Context{Context: context.Background()}, tc.id)
				assert.Equal(t, tc.expectedResponse, team)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestTeamStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.Team(Context{Context: context.Background()}, 123)
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/teams/123",
		Method: http.MethodGet,
	}, err)
}
