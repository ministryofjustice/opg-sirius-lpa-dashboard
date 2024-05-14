package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestUserByEmail(t *testing.T) {
	_, _ := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users"),
						Query: dsl.MapMatcher{
							"email": dsl.String(PotUserEmail),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"id": dsl.Like(47),
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

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
