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

func TestTeam(t *testing.T) {
	pact, err := newPact()

	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/teams/66"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"id":          matchers.Like(66),
							"displayName": matchers.Like("Cool Team"),
							"members": matchers.EachLike(map[string]interface{}{
								"id":          matchers.Like(400),
								"displayName": matchers.Like("Carline"),
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

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

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
