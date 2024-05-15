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

func TestUserByEmail(t *testing.T) {
	pact, err := newPact()

	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func()
		email         string
		expectedUser  User
		expectedError error
	}{
		{
			name:  "OK",
			email: PotUserEmail,
			setup: func() {
				pact.
					AddInteraction().
					Given("!Manager user exists").
					UponReceiving("A request to get !Manager's ID").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/users"),
						Query: matchers.MapMatcher{
							"email": matchers.String(PotUserEmail),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"id": matchers.Like(47),
						}),
					})
			},
			expectedUser: User{
				ID: 47,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				user, err := client.UserByEmail(Context{Context: context.Background()}, tc.email)
				assert.Equal(t, tc.expectedUser, user)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestUserByEmailStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.UserByEmail(Context{Context: context.Background()}, "someone@opgtest.com")
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/users?email=someone@opgtest.com",
		Method: http.MethodGet,
	}, err)
}
