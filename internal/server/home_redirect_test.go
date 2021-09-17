package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockHomeRedirectClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
}

func (m *mockHomeRedirectClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func TestGetHomeRedirect(t *testing.T) {
	testCases := map[string]struct {
		roles            []string
		ExpectedRedirect string
	}{
		"empty": {
			roles:            []string{},
			ExpectedRedirect: "/pending-cases",
		},
		"casework-team": {
			roles:            []string{"Manager"},
			ExpectedRedirect: "/pending-cases",
		},
		"card-payment": {
			roles:            []string{"Card Payment User"},
			ExpectedRedirect: "/tasks",
		},
		"card-payment-manager": {
			roles:            []string{"Manager", "Card Payment User"},
			ExpectedRedirect: "/tasks",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockHomeRedirectClient{}
			client.myDetails.data = sirius.MyDetails{
				Roles: tc.roles,
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/path", nil)

			err := homeRedirect(client)(w, r)
			assert.Equal(RedirectError(tc.ExpectedRedirect), err)
		})
	}
}

func TestPostHomeRedirect(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := feedback(nil, nil)(w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
