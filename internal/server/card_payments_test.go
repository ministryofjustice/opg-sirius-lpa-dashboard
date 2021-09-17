package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockCardPaymentsClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
	tasksByAssignee struct {
		count        int
		lastCtx      sirius.Context
		lastId       int
		lastCriteria sirius.Criteria
		data         []sirius.Task
		pagination   *sirius.Pagination
		err          error
	}
}

func (m *mockCardPaymentsClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockCardPaymentsClient) TasksByAssignee(ctx sirius.Context, id int, criteria sirius.Criteria) ([]sirius.Task, *sirius.Pagination, error) {
	m.tasksByAssignee.count += 1
	m.tasksByAssignee.lastCtx = ctx
	m.tasksByAssignee.lastId = id
	m.tasksByAssignee.lastCriteria = criteria

	return m.tasksByAssignee.data, m.tasksByAssignee.pagination, m.tasksByAssignee.err
}

func TestGetCardPayments(t *testing.T) {
	assert := assert.New(t)

	client := &mockCardPaymentsClient{}
	client.myDetails.data = sirius.MyDetails{
		ID: 14,
	}
	client.tasksByAssignee.data = []sirius.Task{{
		ID: 78,
	}}
	client.tasksByAssignee.pagination = &sirius.Pagination{
		TotalItems: 20,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := cardPayments(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.tasksByAssignee.count)
	assert.Equal(getContext(r), client.tasksByAssignee.lastCtx)
	assert.Equal(14, client.tasksByAssignee.lastId)
	assert.Equal(sirius.Criteria{}.Filter("status", "Not started").Sort("dueDate", sirius.Ascending).Sort("name", sirius.Descending), client.tasksByAssignee.lastCriteria)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(cardPaymentsVars{
		Tasks:             client.tasksByAssignee.data,
		HasIncompleteTask: true,
		XSRFToken:         getContext(r).XSRFToken,
	}, template.lastVars)
}

func TestGetCardPaymentsMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockCardPaymentsClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := cardPayments(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(0, client.tasksByAssignee.count)
}

func TestGetCardPaymentsQueryError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockCardPaymentsClient{}
	client.myDetails.data = sirius.MyDetails{
		ID: 14,
	}
	client.tasksByAssignee.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := cardPayments(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(1, client.tasksByAssignee.count)
}

func TestBadMethodCardPayments(t *testing.T) {
	assert := assert.New(t)

	client := &mockCardPaymentsClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := cardPayments(client, template)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.myDetails.count)
	assert.Equal(0, client.tasksByAssignee.count)
	assert.Equal(0, template.count)
}
