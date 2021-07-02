package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTeam(t *testing.T) {
	pact := newPact()
	defer pact.Teardown()

	testCases := []struct {
		id               int
		name             string
		setup            func()
		cookies          []*http.Cookie
		expectedResponse Team
		expectedError    error
	}{
		{
			name: "OK",
			id:   66,
			setup: func() {
				pact.
					AddInteraction().
					Given("LPA team with members exists").
					UponReceiving("A request for an LPA team").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/66"),
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
							"id":          dsl.Like(66),
							"displayName": dsl.Like("Cool Team"),
							"members": dsl.EachLike(map[string]interface{}{
								"id":          dsl.Like(400),
								"displayName": dsl.Like("Carline"),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: Team{
				ID:          66,
				DisplayName: "Cool Team",
				Members: []TeamMember{
					{
						ID:          400,
						DisplayName: "Carline",
					},
				},
			},
		},
		{
			name: "Unauthorized",
			id:   66,
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request for a team without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/66"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedResponse: Team{},
			expectedError:    ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				team, err := client.Team(getContext(tc.cookies), tc.id)
				assert.Equal(t, tc.expectedResponse, team)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestTeamStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.Team(getContext(nil), 123)
	assert.Equal(t, &StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/teams/123",
		Method: http.MethodGet,
	}, err)
}
