package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockRequestNextTaskClient struct {
	requestNextTask struct {
		count   int
		lastCtx sirius.Context
		err     error
	}
}

func (m *mockRequestNextTaskClient) RequestNextTask(ctx sirius.Context) error {
	m.requestNextTask.count += 1
	m.requestNextTask.lastCtx = ctx

	return m.requestNextTask.err
}

func TestPostRequestNextTask(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequestNextTaskClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := requestNextTask(client)(w, r)
	assert.Equal(RedirectError("/tasks-dashboard"), err)

	assert.Equal(1, client.requestNextTask.count)
	assert.Equal(getContext(r), client.requestNextTask.lastCtx)
}

func TestPostRequestNextTaskErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequestNextTaskClient{}
	client.requestNextTask.err = errors.New("err")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := requestNextTask(client)(w, r)
	assert.Equal(client.requestNextTask.err, err)
}

func TestBadMethodRequestNextTask(t *testing.T) {
	assert := assert.New(t)

	client := &mockRequestNextTaskClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := requestNextTask(client)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.requestNextTask.count)
}
