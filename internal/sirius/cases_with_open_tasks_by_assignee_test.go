package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestCasesWithOpenTasksByAssignee(t *testing.T) {
	pact, err := newPact()
	assert.NoError(t, err)

	testCases := []struct {
		name               string
		setup              func()
		expectedCases      []Case
		expectedPagination *Pagination
		expectedError      error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a case with an open task assigned").
					UponReceiving("A request to get my cases with open tasks").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/assignees/47/cases-with-open-tasks"),
						Query: matchers.MapMatcher{
							"page": matchers.String("1"),
						},
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"total": matchers.Like(1),
							"limit": matchers.Like(25),
							"pages": matchers.Like(map[string]interface{}{
								"current": matchers.Like(1),
								"total":   matchers.Like(1),
							}),
							"cases": matchers.EachLike(map[string]interface{}{
								"id":  matchers.Like(58),
								"uId": matchers.Term("7000-2830-9492", `\d{4}-\d{4}-\d{4}`),
								"donor": matchers.Like(map[string]interface{}{
									"id":        matchers.Like(17),
									"uId":       matchers.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": matchers.Like("Wilma"),
									"surname":   matchers.Like("Ruthman"),
								}),
								"caseSubtype": matchers.Term("hw", "hw|pfa"),
								"receiptDate": matchers.Term("14/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      matchers.String("Pending"),
								"taskCount":   matchers.Like(1),
							}, 1),
						}),
					})
			},
			expectedCases: []Case{{
				ID:  58,
				Uid: "7000-2830-9492",
				Donor: Donor{
					ID:        17,
					Uid:       "7000-5382-4438",
					Firstname: "Wilma",
					Surname:   "Ruthman",
				},
				SubType:     "hw",
				ReceiptDate: SiriusDate{time.Date(2021, 5, 14, 0, 0, 0, 0, time.UTC)},
				Status:      "Pending",
				TaskCount:   1,
			}},
			expectedPagination: &Pagination{
				TotalItems:  1,
				CurrentPage: 1,
				TotalPages:  1,
				PageSize:    25,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				cases, pagination, err := client.CasesWithOpenTasksByAssignee(Context{Context: context.Background()}, 47, 1)
				assert.Equal(t, tc.expectedCases, cases)
				assert.Equal(t, tc.expectedPagination, pagination)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestCasesWithOpenTasksByAssigneeStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, _, err := client.CasesWithOpenTasksByAssignee(Context{Context: context.Background()}, 47, 2)
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/assignees/47/cases-with-open-tasks?page=2",
		Method: http.MethodGet,
	}, err)
}
