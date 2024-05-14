package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	_, _ := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users/47"),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"id":          dsl.Like(47),
							"displayName": dsl.Like("John"),
							"teams": dsl.EachLike(map[string]interface{}{
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

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
