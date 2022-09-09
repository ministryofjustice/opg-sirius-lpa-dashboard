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

func TestTasksByAssignee(t *testing.T) {
	pact := newPact()
	defer pact.Teardown()

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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/47/tasks"),
						Query: dsl.MapMatcher{
							"filter": dsl.String("status:Not started"),
							"sort":   dsl.String("dueDate:asc,name:desc"),
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
							"tasks": dsl.EachLike(map[string]interface{}{
								"id":      dsl.Like(36),
								"status":  dsl.Like("Not started"),
								"dueDate": dsl.Term("19/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"name":    dsl.Like("something"),
								"caseItems": dsl.EachLike(map[string]interface{}{
									"id":  dsl.Like(1),
									"uId": dsl.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
									"donor": dsl.Like(map[string]interface{}{
										"id":        dsl.Like(23),
										"uId":       dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
										"firstname": dsl.Like("Adrian"),
										"surname":   dsl.Like("Kurkjian"),
									}),
									"caseSubtype": dsl.Term("pf", "hw|pf"),
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
					SubType: "pf",
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

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

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
