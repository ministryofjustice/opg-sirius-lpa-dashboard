package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockRequestNextPaymentTaskClient struct {
	requestNextPaymentTask struct {
		count   int
		lastCtx sirius.Context
		err     error
	}
}

func (m *mockRequestNextPaymentTaskClient) RequestNextPaymentTask(ctx sirius.Context) error {
	m.requestNextPaymentTask.count += 1
	m.requestNextPaymentTask.lastCtx = ctx

	return m.requestNextPaymentTask.err
}

func TestPostRequestNextPaymentTask(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequestNextPaymentTaskClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := requestNextPaymentTask(client)(w, r)
	assert.Equal(RedirectError("/card-payments"), err)

	assert.Equal(1, client.requestNextPaymentTask.count)
	assert.Equal(getContext(r), client.requestNextPaymentTask.lastCtx)
}

func TestPostRequestNextPaymentTaskErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequestNextPaymentTaskClient{}
	client.requestNextPaymentTask.err = errors.New("err")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := requestNextPaymentTask(client)(w, r)
	assert.Equal(client.requestNextPaymentTask.err, err)
}

func TestBadMethodRequestNextPaymentTask(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequestNextPaymentTaskClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := requestNextPaymentTask(client)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.requestNextPaymentTask.count)
}
