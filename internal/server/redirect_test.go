package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockRedirectClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
}

func (m *mockRedirectClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func TestRedirect(t *testing.T) {
	assert := assert.New(t)

	client := &mockRedirectClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := redirect(client)(w, r)
	assert.Equal(RedirectError("/pending-cases"), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)
}

func TestRedirectCardPaymentUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockRedirectClient{}
	client.myDetails.data = sirius.MyDetails{Roles: []string{"Card Payment User"}}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := redirect(client)(w, r)
	assert.Equal(RedirectError("/card-payments"), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)
}

func TestRedirectMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockRedirectClient{}
	client.myDetails.err = expectedError

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := redirect(client)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)
}
