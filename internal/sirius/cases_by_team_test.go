package sirius

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestCasesByTeam(t *testing.T) {
	pact := &dsl.Pact{
		Consumer:          "sirius-lpa-dashboard",
		Provider:          "sirius",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		LogDir:            "../../logs",
		PactDir:           "../../pacts",
	}
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
			criteria: Criteria{}.Page(1),
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get my team's cases").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/66/cases"),
						Query: dsl.MapMatcher{
							"page": dsl.String("1"),
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
								"assignee": dsl.Like(map[string]interface{}{
									"id":          dsl.Like(17),
									"displayName": dsl.Like("John Smith"),
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
				Assignee: Assignee{
					ID:          17,
					DisplayName: "John Smith",
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
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("I have a pending case assigned").
					UponReceiving("A request to get my team's cases without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/66/cases"),
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

				cases, pagination, err := client.CasesByTeam(getContext(tc.cookies), 66, tc.criteria)
				assert.Equal(t, tc.expectedCases, cases)
				assert.Equal(t, tc.expectedPagination, pagination)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestCasesByTeamStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, _, err := client.CasesByTeam(getContext(nil), 66, Criteria{}.Page(2))
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/teams/66/cases?page=2",
		Method: http.MethodGet,
	}, err)
}
