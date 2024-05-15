package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestTeams(t *testing.T) {
	pact, err := newPactV2()

	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/teams"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.EachLike(map[string]interface{}{
							"id":          matchers.Like(66),
							"displayName": matchers.Like("Cool Team"),
						}, 1),
					})
			},
			expectedResponse: []Team{
				{
					ID:          66,
					DisplayName: "Cool Team",
					Members:     []TeamMember{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				users, err := client.Teams(Context{Context: context.Background()})
				assert.Equal(t, tc.expectedResponse, users)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestTeamsIgnored(t *testing.T) {
	pact, err := newIgnoredPactV2()

	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/teams"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like([]map[string]interface{}{
							{
								"id":          matchers.Like(66),
								"displayName": matchers.Like("Casework Team"),
								"members": matchers.EachLike(map[string]interface{}{
									"id":          matchers.Like(47),
									"displayName": matchers.Like("John"),
								}, 1),
							},
							{
								"id":          matchers.Like(67),
								"displayName": matchers.Like("Nottingham casework team"),
							},
						}),
					})
			},
			expectedResponse: []Team{
				{
					ID:          66,
					DisplayName: "Casework Team",
					Members: []TeamMember{{
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

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
