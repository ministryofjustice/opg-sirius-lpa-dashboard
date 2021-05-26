package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockRequestNextCasesClient struct {
	requestNextCases struct {
		count   int
		lastCtx sirius.Context
		err     error
	}
}

func (m *mockRequestNextCasesClient) RequestNextCases(ctx sirius.Context) error {
	m.requestNextCases.count += 1
	m.requestNextCases.lastCtx = ctx

	return m.requestNextCases.err
}

func TestPostRequestNextCases(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequestNextCasesClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := requestNextCases(client)(w, r)
	assert.Equal(RedirectError("/pending-cases"), err)

	assert.Equal(1, client.requestNextCases.count)
	assert.Equal(getContext(r), client.requestNextCases.lastCtx)
}

func TestPostRequestNextCasesErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequestNextCasesClient{}
	client.requestNextCases.err = errors.New("err")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := requestNextCases(client)(w, r)
	assert.Equal(client.requestNextCases.err, err)
}

func TestBadMethodRequestNextCases(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequestNextCasesClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := requestNextCases(client)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.requestNextCases.count)
}
