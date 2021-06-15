package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-lpa-dashboard/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockUserTasksClient struct {
	myDetails struct {
		count   int
		lastCtx sirius.Context
		data    sirius.MyDetails
		err     error
	}
	user struct {
		count   int
		lastCtx sirius.Context
		lastId  int
		data    sirius.Assignee
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

func (m *mockUserTasksClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.myDetails.count += 1
	m.myDetails.lastCtx = ctx

	return m.myDetails.data, m.myDetails.err
}

func (m *mockUserTasksClient) User(ctx sirius.Context, id int) (sirius.Assignee, error) {
	m.user.count += 1
	m.user.lastCtx = ctx
	m.user.lastId = id

	return m.user.data, m.user.err
}

func (m *mockUserTasksClient) CasesWithOpenTasksByAssignee(ctx sirius.Context, id, page int) ([]sirius.Case, *sirius.Pagination, error) {
	m.casesWithOpenTasksByAssignee.count += 1
	m.casesWithOpenTasksByAssignee.lastCtx = ctx
	m.casesWithOpenTasksByAssignee.lastId = id
	m.casesWithOpenTasksByAssignee.lastPage = page

	return m.casesWithOpenTasksByAssignee.data, m.casesWithOpenTasksByAssignee.pagination, m.casesWithOpenTasksByAssignee.err
}

func TestGetUserTasks(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserTasksClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.user.data = sirius.Assignee{
		ID:          74,
		DisplayName: "Elfriede Giesing",
	}
	client.casesWithOpenTasksByAssignee.data = []sirius.Case{{
		ID: 78,
		Donor: sirius.Donor{
			ID: 79,
		},
	}}
	client.casesWithOpenTasksByAssignee.pagination = &sirius.Pagination{
		TotalItems: 20,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/tasks/74", nil)

	err := userTasks(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(1, client.casesWithOpenTasksByAssignee.count)
	assert.Equal(getContext(r), client.casesWithOpenTasksByAssignee.lastCtx)
	assert.Equal(74, client.casesWithOpenTasksByAssignee.lastId)
	assert.Equal(1, client.casesWithOpenTasksByAssignee.lastPage)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(userTasksVars{
		Assignee:   client.user.data,
		Cases:      client.casesWithOpenTasksByAssignee.data,
		Pagination: newPagination(client.casesWithOpenTasksByAssignee.pagination),
	}, template.lastVars)
}

func TestGetUserTasksPage(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserTasksClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.user.data = sirius.Assignee{
		ID:          74,
		DisplayName: "Elfriede Giesing",
		Teams: []sirius.Team{{
			ID:          281,
			DisplayName: "Casework Team 6",
		}},
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
	r, _ := http.NewRequest("GET", "/users/tasks/74?page=4", nil)

	err := userTasks(client, template)(w, r)
	assert.Nil(err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(1, client.casesWithOpenTasksByAssignee.count)
	assert.Equal(getContext(r), client.casesWithOpenTasksByAssignee.lastCtx)
	assert.Equal(74, client.casesWithOpenTasksByAssignee.lastId)
	assert.Equal(4, client.casesWithOpenTasksByAssignee.lastPage)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(userTasksVars{
		Assignee:   client.user.data,
		Team:       client.user.data.Teams[0],
		Cases:      client.casesWithOpenTasksByAssignee.data,
		Pagination: newPagination(client.casesWithOpenTasksByAssignee.pagination),
	}, template.lastVars)
}

func TestGetUserTasksForbidden(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserTasksClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Case Worker"},
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/tasks/74", nil)

	err := userTasks(client, template)(w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.user.count)
	assert.Equal(0, client.casesWithOpenTasksByAssignee.count)
}

func TestGetUserTasksMyDetailsError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockUserTasksClient{}
	client.myDetails.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/tasks/74", nil)

	err := userTasks(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(0, client.user.count)
	assert.Equal(0, client.casesWithOpenTasksByAssignee.count)
}

func TestGetUserTasksGetUserError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockUserTasksClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.user.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/tasks/74", nil)

	err := userTasks(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(0, client.casesWithOpenTasksByAssignee.count)
}

func TestGetUserTasksQueryError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockUserTasksClient{}
	client.myDetails.data = sirius.MyDetails{
		Roles: []string{"Manager"},
	}
	client.user.data = sirius.Assignee{
		ID:          74,
		DisplayName: "Elfriede Giesing",
	}
	client.casesWithOpenTasksByAssignee.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/users/tasks/74", nil)

	err := userTasks(client, template)(w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.myDetails.count)
	assert.Equal(getContext(r), client.myDetails.lastCtx)

	assert.Equal(1, client.user.count)
	assert.Equal(getContext(r), client.user.lastCtx)
	assert.Equal(74, client.user.lastId)

	assert.Equal(1, client.casesWithOpenTasksByAssignee.count)
	assert.Equal(getContext(r), client.casesWithOpenTasksByAssignee.lastCtx)
	assert.Equal(74, client.casesWithOpenTasksByAssignee.lastId)
	assert.Equal(1, client.casesWithOpenTasksByAssignee.lastPage)
}

func TestBadMethodUserTasks(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserTasksClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/users/tasks/74", nil)

	err := userTasks(client, template)(w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, client.user.count)
	assert.Equal(0, client.casesWithOpenTasksByAssignee.count)
	assert.Equal(0, template.count)
}
