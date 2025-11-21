package sirius

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestMyDetails(t *testing.T) {
	pact, err := newPact()

	assert.NoError(t, err)

	testCases := []struct {
		name              string
		setup             func()
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
					WithCompleteRequest(consumer.Request{
						Method: http.MethodGet,
						Path:   matchers.String("/lpa-api/v1/users/current"),
					}).
					WithCompleteResponse(consumer.Response{
						Status:  http.StatusOK,
						Headers: matchers.MapMatcher{"Content-Type": matchers.String("application/json")},
						Body: matchers.Like(map[string]interface{}{
							"id":    matchers.Like(47),
							"roles": matchers.EachLike("Manager", 1),
							"teams": matchers.EachLike(map[string]interface{}{
								"id":          matchers.Like(66),
								"displayName": matchers.Like("my team"),
							}, 1),
						}),
					})
			},
			expectedMyDetails: MyDetails{
				ID:    47,
				Roles: []string{"Manager"},
				Teams: []MyDetailsTeam{{ID: 66, DisplayName: "my team"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.ExecuteTest(t, func(config consumer.MockServerConfig) error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://127.0.0.1:%d", config.Port))

				myDetails, err := client.MyDetails(Context{Context: context.Background()})
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

	_, err := client.MyDetails(Context{Context: context.Background()})
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/lpa-api/v1/users/current",
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

func TestMyDetailsHasRole(t *testing.T) {
	testCases := []struct {
		roles    []string
		search   string
		expected bool
	}{
		{
			roles:    []string{},
			search:   "POA User",
			expected: false,
		},
		{
			roles:    []string{"POA User"},
			search:   "POA User",
			expected: true,
		},
		{
			roles:    []string{"Manager", "POA User", "OPG User"},
			search:   "POA User",
			expected: true,
		},
		{
			roles:    []string{"Manager", "OPG User"},
			search:   "POA User",
			expected: false,
		},
	}

	assert := assert.New(t)

	for _, tc := range testCases {
		myDetails := MyDetails{
			Roles: tc.roles,
		}

		assert.Equal(tc.expected, myDetails.HasRole(tc.search))
	}
}

func TestMyDetailsIsSelfAllocationTaskUser(t *testing.T) {
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
			roles:    []string{"Self Allocation Task User"},
			expected: true,
		},
		{
			roles:    []string{"User", "Self Allocation Task User"},
			expected: true,
		},
		{
			roles:    []string{"User", "Self Allocation Task User", "Admin"},
			expected: true,
		},
	}

	assert := assert.New(t)

	for _, tc := range testCases {
		myDetails := MyDetails{
			Roles: tc.roles,
		}

		assert.Equal(tc.expected, myDetails.IsSelfAllocationTaskUser())
	}
}
