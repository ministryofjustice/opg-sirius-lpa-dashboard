package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCasesWithOpenTasksByAssignee(t *testing.T) {
	_, _ := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/47/cases-with-open-tasks"),
						Query: dsl.MapMatcher{
							"page": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"total": dsl.Like(1),
							"limit": dsl.Like(25),
							"pages": dsl.Like(map[string]interface{}{
								"current": dsl.Like(1),
								"total":   dsl.Like(1),
							}),
							"cases": dsl.EachLike(map[string]interface{}{
								"id":  dsl.Like(58),
								"uId": dsl.Term("7000-2830-9492", `\d{4}-\d{4}-\d{4}`),
								"donor": dsl.Like(map[string]interface{}{
									"id":        dsl.Like(17),
									"uId":       dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": dsl.Like("Wilma"),
									"surname":   dsl.Like("Ruthman"),
								}),
								"caseSubtype": dsl.Term("hw", "hw|pfa"),
								"receiptDate": dsl.Term("14/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      dsl.String("Pending"),
								"taskCount":   dsl.Like(1),
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

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
