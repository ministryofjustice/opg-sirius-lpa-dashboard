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

func TestCasesByAssignee(t *testing.T) {
	pact, _ := newPactV2()
	//defer pact.Teardown()

	testCases := []struct {
		name               string
		criteria           Criteria
		setup              func()
		expectedCases      []Case
		expectedPagination *Pagination
		expectedError      error
	}{
		{
			name:     "OK",
			criteria: Criteria{}.Page(1).Sort("receiptDate", Ascending),
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get my cases").
					WithCompleteRequest(
						consumer.Request{
							Method: http.MethodGet,
							Path:   matchers.String("/api/v1/assignees/47/cases"),
							Query: matchers.MapMatcher{
								"page":   matchers.String("1"),
								"filter": matchers.String("caseType:lpa,active:true"),
								"sort":   matchers.String("receiptDate:asc"),
							},
						},
					).
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
								"id":  matchers.Like(36),
								"uId": matchers.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
								"donor": matchers.Like(map[string]interface{}{
									"id":        matchers.Like(23),
									"uId":       matchers.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": matchers.Like("Adrian"),
									"surname":   matchers.Like("Kurkjian"),
								}),
								"caseSubtype": matchers.Term("pfa", "hw|pfa"),
								"receiptDate": matchers.Term("12/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      matchers.Like("Perfect"),
							}, 1),
						}),
					})
			},
			expectedCases: []Case{{
				ID:  36,
				Uid: "7000-8548-8461",
				Donor: Donor{
					ID:        23,
					Uid:       "7000-5382-4438",
					Firstname: "Adrian",
					Surname:   "Kurkjian",
				},
				SubType:     "pfa",
				ReceiptDate: SiriusDate{time.Date(2021, 5, 12, 0, 0, 0, 0, time.UTC)},
				Status:      "Perfect",
			}},
			expectedPagination: &Pagination{
				TotalItems:  1,
				CurrentPage: 1,
				TotalPages:  1,
				PageSize:    25,
			},
		},
		{
			name:     "OK by status",
			criteria: Criteria{}.Filter("status", "Pending").Page(1).Sort("receiptDate", Ascending),
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get my pending cases").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/assignees/47/cases"),
						Query: matchers.MapMatcher{
							"page":   matchers.String("1"),
							"filter": matchers.String("status:Pending,caseType:lpa,active:true"),
							"sort":   matchers.String("receiptDate:asc"),
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
								"status":      matchers.Like("Pending"),
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
			}},
			expectedPagination: &Pagination{
				TotalItems:  1,
				CurrentPage: 1,
				TotalPages:  1,
				PageSize:    25,
			},
		},
		{
			name:     "OK by workedDate",
			criteria: Criteria{}.Filter("status", "Pending").Page(1).Sort("workedDate", Descending).Sort("receiptDate", Ascending),
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get my pending cases sorted by worked date").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/assignees/47/cases"),
						Query: matchers.MapMatcher{
							"page":   matchers.String("1"),
							"filter": matchers.String("status:Pending,caseType:lpa,active:true"),
							"sort":   matchers.String("workedDate:desc,receiptDate:asc"),
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
								"status":      matchers.Like("Pending"),
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
			}},
			expectedPagination: &Pagination{
				TotalItems:  1,
				CurrentPage: 1,
				TotalPages:  1,
				PageSize:    25,
			},
		},
		{
			name:     "OK with criteria",
			criteria: Criteria{}.Filter("status", "Pending").Page(1).Limit(1).Sort("receiptDate", Ascending),
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get my oldest pending case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/assignees/47/cases"),
						Query: matchers.MapMatcher{
							"page":   matchers.String("1"),
							"limit":  matchers.String("1"),
							"filter": matchers.String("status:Pending,caseType:lpa,active:true"),
							"sort":   matchers.String("receiptDate:asc"),
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
								"id":  matchers.Like(453),
								"uId": matchers.Term("7000-2830-9429", `\d{4}-\d{4}-\d{4}`),
								"donor": matchers.Like(map[string]interface{}{
									"id":        matchers.Like(363),
									"uId":       matchers.Term("7000-5382-4435", `\d{4}-\d{4}-\d{4}`),
									"firstname": matchers.Like("Mario"),
									"surname":   matchers.Like("Evanosky"),
								}),
								"caseSubtype": matchers.Term("hw", "hw|pfa"),
								"receiptDate": matchers.Term("28/11/2017", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      matchers.Like("Pending"),
							}, 1),
						}),
					})
			},
			expectedCases: []Case{{
				ID:  453,
				Uid: "7000-2830-9429",
				Donor: Donor{
					ID:        363,
					Uid:       "7000-5382-4435",
					Firstname: "Mario",
					Surname:   "Evanosky",
				},
				SubType:     "hw",
				ReceiptDate: SiriusDate{time.Date(2017, 11, 28, 0, 0, 0, 0, time.UTC)},
				Status:      "Pending",
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

				cases, pagination, err := client.CasesByAssignee(Context{Context: context.Background()}, 47, tc.criteria)
				assert.Equal(t, tc.expectedCases, cases)
				assert.Equal(t, tc.expectedPagination, pagination)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestHasWorkableCase(t *testing.T) {
	pact, _ := newPactV2()

	testCases := []struct {
		name          string
		criteria      Criteria
		setup         func()
		expectedError error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get a pending, unworked case").
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/api/v1/assignees/47/cases"),
						Query: matchers.MapMatcher{
							"page":   matchers.String("1"),
							"limit":  matchers.String("1"),
							"filter": matchers.String("status:Pending,worked:false,caseType:lpa,active:true"),
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
								"id":  matchers.Like(453),
								"uId": matchers.Term("7000-2830-9429", `\d{4}-\d{4}-\d{4}`),
								"donor": matchers.Like(map[string]interface{}{
									"id":        matchers.Like(363),
									"uId":       matchers.Term("7000-5382-4435", `\d{4}-\d{4}-\d{4}`),
									"firstname": matchers.Like("Mario"),
									"surname":   matchers.Like("Evanosky"),
								}),
								"caseSubtype": matchers.Term("hw", "hw|pfa"),
								"receiptDate": matchers.Term("28/11/2017", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      matchers.Like("Pending"),
							}, 1),
						}),
					})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				ok, err := client.HasWorkableCase(Context{Context: context.Background()}, 47)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestCasesByAssigneeStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, _, err := client.CasesByAssignee(Context{Context: context.Background()}, 47, Criteria{}.Page(2))
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/assignees/47/cases?filter=caseType%3Alpa%2Cactive%3Atrue&page=2",
		Method: http.MethodGet,
	}, err)
}

func TestHasWorkableCaseStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.HasWorkableCase(Context{Context: context.Background()}, 47)
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/assignees/47/cases?filter=status%3APending%2Cworked%3Afalse%2CcaseType%3Alpa%2Cactive%3Atrue&limit=1&page=1",
		Method: http.MethodGet,
	}, err)
}
