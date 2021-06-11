package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestMyDetails(t *testing.T) {
	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		name              string
		setup             func()
		cookies           []*http.Cookie
		expectedMyDetails MyDetails
		expectedError     error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get my details").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users/current"),
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
							"id":    dsl.Like(47),
							"roles": dsl.EachLike("Manager", 1),
							"teams": dsl.EachLike(map[string]interface{}{
								"id":          dsl.Like(66),
								"displayName": dsl.Like("my team"),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedMyDetails: MyDetails{
				ID:    47,
				Roles: []string{"Manager"},
				Teams: []MyDetailsTeam{{ID: 66, DisplayName: "my team"}},
			},
		},

		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get my details without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/users/current"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
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

				myDetails, err := client.MyDetails(getContext(tc.cookies))
				assert.Equal(t, tc.expectedMyDetails, myDetails)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestMyDetailsStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.MyDetails(getContext(nil))
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/users/current",
		Method: http.MethodGet,
	}, err)
}

func TestMyDetailsIsManager(t *testing.T) {
	testCases := []struct {
		roles    []string
		expected bool
	}{
		{
			roles:    []string{},
			expected: false,
		},
		{
			roles:    []string{"Admin"},
			expected: false,
		},
		{
			roles:    []string{"Manager"},
			expected: true,
		},
		{
			roles:    []string{"User", "Admin"},
			expected: false,
		},
		{
			roles:    []string{"User", "Manager", "Admin"},
			expected: true,
		},
	}

	assert := assert.New(t)

	for _, tc := range testCases {
		myDetails := MyDetails{
			Roles: tc.roles,
		}

		assert.Equal(tc.expected, myDetails.IsManager())
	}
}

func TestMyDetailsIsCaseWorker(t *testing.T) {
	testCases := []struct {
		roles    []string
		expected bool
	}{
		{
			roles:    []string{},
			expected: false,
		},
		{
			roles:    []string{"OPG User"},
			expected: false,
		},
		{
			roles:    []string{"OPG User", "POA User"},
			expected: true,
		},
		{
			roles:    []string{"POA User", "OPG User"},
			expected: true,
		},
		{
			roles:    []string{"OPG User", "POA User", "Manager"},
			expected: false,
		},
	}

	assert := assert.New(t)

	for _, tc := range testCases {
		myDetails := MyDetails{
			Roles: tc.roles,
		}

		assert.Equal(tc.expected, myDetails.IsCaseWorker())
	}
}
