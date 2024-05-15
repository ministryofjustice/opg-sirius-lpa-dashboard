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

func TestUser(t *testing.T) {
	pact, err := newPactV2()

	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		email         string
		expectedUser  Assignee
		expectedError error
	}{
		{
			name:  "OK",
			email: PotUserEmail,
			setup: func() {
				pact.
					AddInteraction().
					Given("!Manager user exists").
					UponReceiving("A request to get a user").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/users/47"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"id":          matchers.Like(47),
							"displayName": matchers.Like("John"),
							"teams": matchers.EachLike(map[string]interface{}{
								"id":          66,
								"displayName": "Cool Team",
							}, 1),
						}),
					})
			},
			expectedUser: Assignee{
				ID:          47,
				DisplayName: "John",
				Teams: []Team{
					{
						ID:          66,
						DisplayName: "Cool Team",
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

				user, err := client.User(Context{Context: context.Background()}, 47)
				assert.Equal(t, tc.expectedUser, user)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestUserStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.User(Context{Context: context.Background()}, 47)
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/users/47",
		Method: http.MethodGet,
	}, err)
}
