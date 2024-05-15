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

func TestTasksByAssignee(t *testing.T) {
	pact, err := newPactV2()

	assert.NoError(t, err)

	testCases := []struct {
		name               string
		criteria           Criteria
		setup              func()
		expectedTasks      []Task
		expectedPagination *Pagination
		expectedError      error
	}{
		{
			name:     "OK",
			criteria: Criteria{}.Filter("status", "Not started").Sort("dueDate", Ascending).Sort("name", Descending),
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a task assigned").
					UponReceiving("A request to get my tasks").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/assignees/47/tasks"),
						Query: matchers.MapMatcher{
							"filter": matchers.String("status:Not started"),
							"sort":   matchers.String("dueDate:asc,name:desc"),
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
							"tasks": matchers.EachLike(map[string]interface{}{
								"id":      matchers.Like(36),
								"status":  matchers.Like("Not started"),
								"dueDate": matchers.Term("19/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"name":    matchers.Like("something"),
								"caseItems": matchers.EachLike(map[string]interface{}{
									"id":  matchers.Like(1),
									"uId": matchers.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
									"donor": matchers.Like(map[string]interface{}{
										"id":        matchers.Like(23),
										"uId":       matchers.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
										"firstname": matchers.Like("Adrian"),
										"surname":   matchers.Like("Kurkjian"),
									}),
									"caseSubtype": matchers.Term("pfa", "hw|pfa"),
								}, 1),
							}, 1),
						}),
					})
			},
			expectedTasks: []Task{{
				ID:      36,
				Status:  "Not started",
				DueDate: SiriusDate{time.Date(2021, 5, 19, 0, 0, 0, 0, time.UTC)},
				Name:    "something",
				CaseItems: []TaskCaseItem{{
					ID:  1,
					Uid: "7000-8548-8461",
					Donor: Donor{
						ID:        23,
						Uid:       "7000-5382-4438",
						Firstname: "Adrian",
						Surname:   "Kurkjian",
					},
					SubType: "pfa",
				}},
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

				tasks, pagination, err := client.TasksByAssignee(Context{Context: context.Background()}, 47, tc.criteria)
				assert.Equal(t, tc.expectedTasks, tasks)
				assert.Equal(t, tc.expectedPagination, pagination)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestTasksByAssigneeStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, _, err := client.TasksByAssignee(Context{Context: context.Background()}, 47, Criteria{})
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/assignees/47/tasks?",
		Method: http.MethodGet,
	}, err)
}
