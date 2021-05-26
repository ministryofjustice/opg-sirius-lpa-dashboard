package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestUserByEmail(t *testing.T) {
	pact := &dsl.Pact{
		Consumer:          "sirius-lpa-dashboard",
		Provider:          "sirius",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		LogDir:            "../../logs",
		PactDir:           "../../pacts",
	}
	defer pact.Teardown()

	testCases := []struct {
		name          string
		setup         func()
		cookies       []*http.Cookie
		email         string
		expectedUser  User
		expectedError error
	}{
		{
			name:  "OK",
			email: "manager@opgtest.com",
			setup: func() {
				pact.
					AddInteraction().
					Given("!Manager user exists").
					UponReceiving("A request to get !Manager's ID").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users"),
						Query: dsl.MapMatcher{
							"email": dsl.String("manager@opgtest.com"),
						},
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
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
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedUser: User{
				ID: 47,
			},
		},

		{
			name:  "Unauthorized",
			email: "manager@opgtest.com",
			setup: func() {
				pact.
					AddInteraction().
					Given("!Manager user exists").
					UponReceiving("A request to get !Manager's ID without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users"),
						Query: dsl.MapMatcher{
							"email": dsl.String("manager@opgtest.com"),
						},
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				user, err := client.UserByEmail(getContext(tc.cookies), tc.email)
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

	_, err := client.UserByEmail(getContext(nil), "someone@opgtest.com")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/users?email=someone@opgtest.com",
		Method: http.MethodGet,
	}, err)
}
