package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockMarkWorkedClient struct {
	markWorked struct {
		count   int
		lastCtx sirius.Context
		err     error
		ids     []int
	}
}

func (m *mockMarkWorkedClient) MarkWorked(ctx sirius.Context, id int) error {
	m.markWorked.count += 1
	m.markWorked.lastCtx = ctx
	m.markWorked.ids = append(m.markWorked.ids, id)

	return m.markWorked.err
}

func TestPostMarkWorked(t *testing.T) {
	assert := assert.New(t)

	client := &mockMarkWorkedClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("worked=12&worked=34"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := markWorked(client)(w, r)
	assert.Equal(RedirectError("/pending-cases"), err)

	assert.Equal(2, client.markWorked.count)
	assert.Equal(getContext(r), client.markWorked.lastCtx)
	assert.Equal([]int{12, 34}, client.markWorked.ids)
}

func TestPostMarkWorkedNoForm(t *testing.T) {
	assert := assert.New(t)

	client := &mockMarkWorkedClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := markWorked(client)(w, r)
	assert.NotNil(err)

	assert.Equal(0, client.markWorked.count)
}

func TestPostMarkWorkedBadId(t *testing.T) {
	assert := assert.New(t)

	client := &mockMarkWorkedClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("worked=what"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := markWorked(client)(w, r)
	assert.NotNil(err)

	assert.Equal(0, client.markWorked.count)
}

func TestPostMarkWorkedErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockMarkWorkedClient{}
	client.markWorked.err = errors.New("err")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("worked=1"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := markWorked(client)(w, r)
	assert.Equal(client.markWorked.err, err)
}

func TestBadMethodMarkWorked(t *testing.T) {
	assert := assert.New(t)

	client := &mockMarkWorkedClient{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := markWorked(client)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.markWorked.count)
}
