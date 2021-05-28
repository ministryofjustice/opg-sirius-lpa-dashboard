package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockLandingPageClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
}

func (m *mockLandingPageClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func TestGetLandingPageManager(t *testing.T) {
	assert := assert.New(t)

	client := &mockLandingPageClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := landingPage(client)(w, r)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(RedirectError("/teams/central"), err)
}

func TestGetLandingPageCaseWorker(t *testing.T) {
	assert := assert.New(t)

	client := &mockLandingPageClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := landingPage(client)(w, r)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(RedirectError("/pending-cases"), err)
}

func TestGetLandingPageMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockLandingPageClient{}
	client.myDetails.err = expectedError

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := landingPage(client)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)
}

func TestBadMethodLandingPage(t *testing.T) {
	assert := assert.New(t)

	client := &mockLandingPageClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := landingPage(client)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
