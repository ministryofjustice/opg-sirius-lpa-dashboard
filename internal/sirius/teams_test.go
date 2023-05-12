package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTeams(t *testing.T) {
	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []Team
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists and teams have no type").
					UponReceiving("A request for teams").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like([]map[string]interface{}{
							{
								"id":          dsl.Like(66),
								"displayName": dsl.Like("Casework Team"),
							},
							{
								"id":          dsl.Like(67),
								"displayName": dsl.Like("Nottingham casework team"),
							},
						}),
					})
			},
			expectedResponse: []Team{
				{
					ID:          66,
					DisplayName: "Casework Team",
					Members:     []TeamMember{},
				},
				{
					ID:          67,
					DisplayName: "Nottingham casework team",
					Members:     []TeamMember{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				users, err := client.Teams(Context{Context: context.Background()})
				assert.Equal(t, tc.expectedResponse, users)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestTeamsIgnored(t *testing.T) {
	pact := newIgnoredPact()
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		expectedResponse []Team
		expectedError    error
	}{
		{
			name: "OK with members",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists and teams have no type").
					UponReceiving("A request for teams").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like([]map[string]interface{}{
							{
								"id":          dsl.Like(66),
								"displayName": dsl.Like("Casework Team"),
								"members": dsl.EachLike(map[string]interface{}{
									"id":          dsl.Like(47),
									"displayName": dsl.Like("John"),
								}, 1),
							},
							{
								"id":          dsl.Like(67),
								"displayName": dsl.Like("Nottingham casework team"),
							},
						}),
					})
			},
			expectedResponse: []Team{
				{
					ID:          66,
					DisplayName: "Casework Team",
					Members:     []TeamMember{{
						ID:          47,
						DisplayName: "John",
					}},
				},
				{
					ID:          67,
					DisplayName: "Nottingham casework team",
					Members:     []TeamMember{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				users, err := client.Teams(Context{Context: context.Background()})
				assert.Equal(t, tc.expectedResponse, users)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestTeamsStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.Teams(Context{Context: context.Background()})
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/teams",
		Method: http.MethodGet,
	}, err)
}
