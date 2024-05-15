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

func TestFeedback(t *testing.T) {
	pact, err := newIgnoredPactV2()

	assert.NoError(t, err)

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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPost,
						Path:   matchers.String("/api/wth"),
						Body: matchers.Like(map[string]interface{}{
							"message": matchers.String("hey"),
						}),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

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
