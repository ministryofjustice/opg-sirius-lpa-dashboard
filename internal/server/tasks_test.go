package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockTasksClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
	casesWithOpenTasksByAssignee struct {
		count      int
		lastCtx    sirius.Context
		lastId     int
		lastPage   int
		data       []sirius.Case
		pagination *sirius.Pagination
		err        error
	}
}

func (m *mockTasksClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockTasksClient) CasesWithOpenTasksByAssignee(ctx sirius.Context, id, page int) ([]sirius.Case, *sirius.Pagination, error) {
	m.casesWithOpenTasksByAssignee.count += 1
	m.casesWithOpenTasksByAssignee.lastCtx = ctx
	m.casesWithOpenTasksByAssignee.lastId = id
	m.casesWithOpenTasksByAssignee.lastPage = page

	return m.casesWithOpenTasksByAssignee.data, m.casesWithOpenTasksByAssignee.pagination, m.casesWithOpenTasksByAssignee.err
}

func TestGetTasks(t *testing.T) {
	testCases := map[string]struct {
		URL  string
		Page int
	}{
		"no-page": {
			URL:  "/path",
			Page: 1,
		},
		"specified-page": {
			URL:  "/path?page=5",
			Page: 5,
		},
		"bad-page": {
			URL:  "/path?page=what",
			Page: 1,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockTasksClient{}
			client.myDetails.data = sirius.MyDetails{
				ID: 14,
			}
			client.casesWithOpenTasksByAssignee.data = []sirius.Case{{
				ID: 78,
				Donor: sirius.Donor{
					ID: 79,
				},
			}}
			client.casesWithOpenTasksByAssignee.pagination = &sirius.Pagination{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", tc.URL, nil)

			err := tasks(client, template)(w, r)
			assert.Nil(err)

			assert.Equal(1, client.myDetails.count)
			assert.Equal(getContext(r), client.myDetails.lastCtx)

			assert.Equal(1, client.casesWithOpenTasksByAssignee.count)
			assert.Equal(getContext(r), client.casesWithOpenTasksByAssignee.lastCtx)
			assert.Equal(14, client.casesWithOpenTasksByAssignee.lastId)
			assert.Equal(tc.Page, client.casesWithOpenTasksByAssignee.lastPage)

			assert.Equal(1, template.count)
			assert.Equal("page", template.lastName)
			assert.Equal(tasksVars{
				Cases:      client.casesWithOpenTasksByAssignee.data,
				Pagination: client.casesWithOpenTasksByAssignee.pagination,
			}, template.lastVars)
		})
	}
}

func TestGetTasksMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockTasksClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := tasks(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.casesWithOpenTasksByAssignee.count)
}

func TestGetTasksQueryError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockTasksClient{}
	client.myDetails.data = sirius.MyDetails{
		ID: 14,
	}
	client.casesWithOpenTasksByAssignee.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := tasks(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.casesWithOpenTasksByAssignee.count)
	assert.Equal(getContext(r), client.casesWithOpenTasksByAssignee.lastCtx)
	assert.Equal(14, client.casesWithOpenTasksByAssignee.lastId)
}

func TestBadMethodTasks(t *testing.T) {
	assert := assert.New(t)

	client := &mockTasksClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/path", nil)

	err := tasks(client, template)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.myDetails.count)
	assert.Equal(0, client.casesWithOpenTasksByAssignee.count)
	assert.Equal(0, template.count)
}
