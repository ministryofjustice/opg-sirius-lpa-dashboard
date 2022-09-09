package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestAssign(t *testing.T) {
	pact := newPact()
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
					Given("I have a pending case assigned").
					UponReceiving("A request to reassign a case").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/users/47/cases/58"),
						Body: dsl.Like(map[string]interface{}{
							"data": dsl.EachLike(map[string]interface{}{
								"assigneeId": dsl.Like(99),
								"caseType":   dsl.String("LPA"),
								"id":         dsl.Like(1),
							}, 1),
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
