package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestFeedback(t *testing.T) {
	pact := newIgnoredPact()
	defer pact.Teardown()

	testCases := []struct {
		name          string
		setup         func()
		expectedError error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a user").
					UponReceiving("A request to give feedback").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/api/wth"),
						Body: dsl.Like(map[string]interface{}{
							"message": dsl.String("hey"),
						}),
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.Feedback(Context{Context: context.Background()}, "hey")
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestFeedbackStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.Feedback(Context{Context: context.Background()}, "hey")
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/wth",
		Method: http.MethodPost,
	}, err)
}
