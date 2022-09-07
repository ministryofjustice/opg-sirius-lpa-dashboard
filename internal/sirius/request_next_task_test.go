package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestRequestNextTask(t *testing.T) {
	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name          string
		setup         func()
		cookies       []*http.Cookie
		expectedError error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have no assigned tasks and there is an available payment task").
					UponReceiving("A request to be assigned a new payment task").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/request-new-task"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
		},
		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have no assigned tasks and there is an available payment task").
					UponReceiving("A request to be assigned new payment task without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/v1/request-new-task"),
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

				err := client.RequestNextTask(getContext(tc.cookies))
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestRequestNextTaskStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.RequestNextTask(getContext(nil))
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/request-new-task",
		Method: http.MethodPost,
	}, err)
}
