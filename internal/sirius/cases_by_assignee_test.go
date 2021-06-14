package sirius

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCasesByAssignee(t *testing.T) {
	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name               string
		criteria           Criteria
		setup              func()
		cookies            []*http.Cookie
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/47/cases"),
						Query: dsl.MapMatcher{
							"page":   dsl.String("1"),
							"filter": dsl.String("caseType:lpa,active:true"),
							"sort":   dsl.String("receiptDate:asc"),
						},
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
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
								"id":  dsl.Like(36),
								"uId": dsl.Term("7000-8548-8461", `\d{4}-\d{4}-\d{4}`),
								"donor": dsl.Like(map[string]interface{}{
									"id":        dsl.Like(23),
									"uId":       dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": dsl.Like("Adrian"),
									"surname":   dsl.Like("Kurkjian"),
								}),
								"caseSubtype": dsl.Term("pf", "hw|pf"),
								"receiptDate": dsl.Term("12/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      dsl.Like("Perfect"),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
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
				SubType:     "pf",
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/47/cases"),
						Query: dsl.MapMatcher{
							"page":   dsl.String("1"),
							"filter": dsl.String("status:Pending,caseType:lpa,active:true"),
							"sort":   dsl.String("receiptDate:asc"),
						},
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
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
								"caseSubtype": dsl.Term("hw", "hw|pf"),
								"receiptDate": dsl.Term("14/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      dsl.Like("Pending"),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
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
			criteria: Criteria{}.Filter("status", "Pending").Page(1).Sort("workedDate", Descending),
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get my pending cases sorted by worked date").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/47/cases"),
						Query: dsl.MapMatcher{
							"page":   dsl.String("1"),
							"filter": dsl.String("status:Pending,caseType:lpa,active:true"),
							"sort":   dsl.String("workedDate:desc"),
						},
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
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
								"caseSubtype": dsl.Term("hw", "hw|pf"),
								"receiptDate": dsl.Term("14/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      dsl.Like("Pending"),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
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
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/47/cases"),
						Query: dsl.MapMatcher{
							"page":   dsl.String("1"),
							"limit":  dsl.String("1"),
							"filter": dsl.String("status:Pending,caseType:lpa,active:true"),
							"sort":   dsl.String("receiptDate:asc"),
						},
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
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
								"id":  dsl.Like(453),
								"uId": dsl.Term("7000-2830-9429", `\d{4}-\d{4}-\d{4}`),
								"donor": dsl.Like(map[string]interface{}{
									"id":        dsl.Like(363),
									"uId":       dsl.Term("7000-5382-4435", `\d{4}-\d{4}-\d{4}`),
									"firstname": dsl.Like("Mario"),
									"surname":   dsl.Like("Evanosky"),
								}),
								"caseSubtype": dsl.Term("hw", "hw|pf"),
								"receiptDate": dsl.Term("28/11/2017", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      dsl.Like("Pending"),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
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
		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get my cases without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/47/cases"),
						Query: dsl.MapMatcher{
							"filter": dsl.String("caseType:lpa,active:true"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				cases, pagination, err := client.CasesByAssignee(getContext(tc.cookies), 47, tc.criteria)
				assert.Equal(t, tc.expectedCases, cases)
				assert.Equal(t, tc.expectedPagination, pagination)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestHasWorkableCase(t *testing.T) {
	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name          string
		criteria      Criteria
		setup         func()
		cookies       []*http.Cookie
		expectedError error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get a pending, unworked case").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/47/cases"),
						Query: dsl.MapMatcher{
							"page":   dsl.String("1"),
							"limit":  dsl.String("1"),
							"filter": dsl.String("status:Pending,worked:false,caseType:lpa,active:true"),
						},
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
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
								"id":  dsl.Like(453),
								"uId": dsl.Term("7000-2830-9429", `\d{4}-\d{4}-\d{4}`),
								"donor": dsl.Like(map[string]interface{}{
									"id":        dsl.Like(363),
									"uId":       dsl.Term("7000-5382-4435", `\d{4}-\d{4}-\d{4}`),
									"firstname": dsl.Like("Mario"),
									"surname":   dsl.Like("Evanosky"),
								}),
								"caseSubtype": dsl.Term("hw", "hw|pf"),
								"receiptDate": dsl.Term("28/11/2017", `\d{1,2}/\d{1,2}/\d{4}`),
								"status":      dsl.Like("Pending"),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				ok, err := client.HasWorkableCase(getContext(tc.cookies), 47)
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

	_, _, err := client.CasesByAssignee(getContext(nil), 47, Criteria{}.Page(2))
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/assignees/47/cases?filter=caseType%3Alpa%2Cactive%3Atrue&page=2",
		Method: http.MethodGet,
	}, err)
}
