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

func TestAssign(t *testing.T) {

	pact, err := newPactV2()

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
					UponReceiving("A request to reassign a case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodPut,
						Path:   matchers.String("/api/v1/users/47/cases/58"),
						Body: matchers.Like(map[string]interface{}{
							"data": matchers.EachLike(map[string]interface{}{
								"assigneeId": matchers.Like(99),
								"caseType":   matchers.String("LPA"),
								"id":         matchers.Like(1),
							}, 1),
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

				err := client.Assign(Context{Context: context.Background()}, []int{58}, 47)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestAssignStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.Assign(Context{Context: context.Background()}, []int{1}, 47)
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/users/47/cases/1",
		Method: http.MethodPut,
	}, err)
}
