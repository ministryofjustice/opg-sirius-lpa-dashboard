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
		name          string
		setup         func()
		cookies       []*http.Cookie
		expectedCases []Case
		expectedError error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("A user has cases").
					UponReceiving("A request to get a user's cases").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/47/cases"),
						Query: dsl.MapMatcher{
							"filter": dsl.String("caseType:lpa,status:Pending,active:true"),
							"sort":   dsl.String("caseSubType:asc"),
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
							"cases": dsl.EachLike(map[string]interface{}{
								"id":  dsl.Like(58),
								"uId": dsl.Term("7000-2830-9492", `\d{4}-\d{4}-\d{4}`),
								"donor": dsl.Like(map[string]interface{}{
									"id":        dsl.Like(17),
									"uId":       dsl.Term("7000-5382-4438", `\d{4}-\d{4}-\d{4}`),
									"firstname": dsl.Like("Wilma"),
									"surname":   dsl.Like("Ruthman"),
								}),
								"caseSubtype": dsl.Term("HW", "HW|PF"),
								"receiptDate": dsl.Term("14/05/2021", `\d{1,2}/\d{1,2}/\d{4}`),
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
					FirstName: "Wilma",
					Surname:   "Ruthman",
				},
				SubType:     "HW",
				ReceiptDate: SiriusDate{time.Date(2021, 5, 14, 0, 0, 0, 0, time.UTC)},
			}},
		},

		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("A user has cases").
					UponReceiving("A request to get a user's cases without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/47/cases"),
						Query: dsl.MapMatcher{
							"filter": dsl.String("caseType:lpa,status:Pending,active:true"),
							"sort":   dsl.String("caseSubType:asc"),
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

				cases, err := client.CasesByAssignee(getContext(tc.cookies), 47)
				assert.Equal(t, tc.expectedCases, cases)
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

	_, err := client.CasesByAssignee(getContext(nil), 47)
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/assignees/47/cases?filter=caseType:lpa,status:Pending,active:true&sort=caseSubType:asc",
		Method: http.MethodGet,
	}, err)
}
