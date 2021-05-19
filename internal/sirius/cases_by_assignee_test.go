package sirius

import (
	"fmt"
	"net/http"
	"testing"

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
						Path:   dsl.String("/api/v1/assignees/123/cases"),
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
								"id": dsl.Like(58),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedCases: []Case{{
				ID: 58,
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
						Path:   dsl.String("/api/v1/assignees/123/cases"),
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

				cases, err := client.CasesByAssignee(getContext(tc.cookies), 123)
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

	_, err := client.CasesByAssignee(getContext(nil), 123)
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/assignees/123/cases?filter=caseType:lpa,status:Pending,active:true&sort=caseSubType:asc",
		Method: http.MethodGet,
	}, err)
}
