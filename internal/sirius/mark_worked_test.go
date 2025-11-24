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

func TestMarkWorked(t *testing.T) {
	pact, err := newPact()
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
					Given("I have a pending case assigned").
					UponReceiving("A request to mark the case as worked").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/lpa-api/v1/lpas/800"),
						Body:   matchers.Like(map[string]interface{}{"worked": true}),
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

				err := client.MarkWorked(Context{Context: context.Background()}, 800)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestMarkWorkedStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.MarkWorked(Context{Context: context.Background()}, 1)
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/lpa-api/v1/lpas/1",
		Method: http.MethodPut,
	}, err)
}
